# goSysMon – A Linux System Monitor CLI Tool

A real‑time terminal system monitor written in **Go**, built **from scratch** using only the Linux `/proc` filesystem.  
No external monitoring libraries – just `os`, `bufio`, and `/proc` parsing.

---

## Features

- Reads **CPU usage** (`/proc/stat`) with delta‑based percentage calculation
- Shows **memory stats** (`/proc/meminfo`): total, available, used
- Lists **running processes** by scanning `/proc/[PID]` directories
- Displays per‑process **PID**, **command**, **CPU‑time ticks**, and **memory (RSS)**
- Auto‑refreshes at configurable intervals (CLI flag)
- Handles **SIGINT** / **SIGTERM** gracefully with a shutdown message
- Structured into idiomatic Go packages: `proc`, `monitor`, `display`
- **Unit tests** for all `/proc` parsing functions

---

## Prerequisites

- **Linux** – the tool reads `/proc`, which only exists on Linux.
  - You can build on macOS but must run on a Linux machine, VM, or container.
- **Go** 1.20+ (latest stable recommended)

---

## Installation & Build

```bash
# Clone the repository
git clone https://github.com/Jephtaah/goSysMon.git
cd goSysMon

# Build the binary
go build -o goSysMon .