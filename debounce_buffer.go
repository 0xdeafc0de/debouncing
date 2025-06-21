package debounce

import (
	"sync"
	"time"
)

type Manager struct {
	cfg     Config
	buffers map[string]*buffer
	mu      sync.Mutex
}

type Config struct {
	FlushAfter                time.Duration
	MinIntervalBetweenFlushes time.Duration
	SendFunc                  func(string, []interface{})
}

type buffer struct {
	mu             sync.Mutex
	items          []interface{}
	timer          *time.Timer
	lastFlushTime  time.Time
	flushScheduled bool
}

func NewManager(cfg Config) *Manager {
	return &Manager{
		cfg:     cfg,
		buffers: make(map[string]*buffer),
	}
}

func (m *Manager) Add(key string, item interface{}) {
	m.mu.Lock()
	b, exists := m.buffers[key]
	if !exists {
		b = &buffer{}
		m.buffers[key] = b
	}
	m.mu.Unlock()

	b.mu.Lock()
	b.items = append(b.items, item)
	if !b.flushScheduled {
		b.flushScheduled = true
		b.timer = time.AfterFunc(m.cfg.FlushAfter, func() {
			m.flush(key)
		})
	}
	b.mu.Unlock()
}

func (m *Manager) flush(key string) {
	m.mu.Lock()
	b, exists := m.buffers[key]
	if !exists {
		m.mu.Unlock()
		return
	}
	m.mu.Unlock()

	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	if !b.lastFlushTime.IsZero() && now.Sub(b.lastFlushTime) < m.cfg.MinIntervalBetweenFlushes {
		delay := m.cfg.MinIntervalBetweenFlushes - now.Sub(b.lastFlushTime)
		b.timer = time.AfterFunc(delay, func() {
			m.flush(key)
		})
		return
	}

	batch := b.items
	b.items = nil
	b.lastFlushTime = now
	b.flushScheduled = false

	go m.cfg.SendFunc(key, batch)
}
