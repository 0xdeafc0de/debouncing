# Debounce
A **generic Go library** for **debouncing** and **batching** signals from multiple sources (e.g., devices, users, sensors), with **per-key isolation** and **rate-limited flushing**.

This package helps solve a common systems problem:

> You receive a bursty stream of signals for many different keys (e.g., device IDs or MAC addresses), but want to **batch and flush** those signals:
> 
> - Within a **maximum duration** (e.g., 2 seconds),
> - With a **minimum interval** between flushes (e.g., 200ms),
> - **Independently for each key**.

This is especially useful when sending data to downstream systems (like a control plane, queue, or database) that shouldn't be overwhelmed by bursts.


---

## Features

- **Debouncing**: Collect items for a short period before sending.
- **Per-Key Isolation**: Independent buffers and timers for each key.
- **Rate-Limited Flushing**: Avoids excessive flushes with `MinIntervalBetweenFlushes`.
- **Custom Send Function**: Define how you want to handle flushed batches.

---
