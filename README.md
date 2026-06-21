# cho

A Log collector, written in Go. `cho` runs as a sidecar alongside your application, tails its log file, parses structured JSON logs, and batches them for delivery to a log ingestor (`calatrava`).

Inspired by the operational layer built around [jkvs](https://github.com/persona-mp3/jkvs) — a persistent key-value store that `cho` was originally built to observe.

---

## How It Works

```
application process
  └── writes structured logs → log file
                                    │
                              cho (sidecar)
                              ├── fsnotify watches directory for writes  [+]
                              ├── ReadAt tracks byte offset between reads
                              ├── buffers parsed Log entries
                              └── flushes when threshold or interval fires
                                    │
                                    ▼
                              calatrava (ingestor)   ← in progress
```

`cho` never touches the log file from the start on every tick — it tracks the last read offset and reads only new bytes since the previous read. Multiple write events from a single logical write are coalesced naturally: the offset advances by however many bytes were written, regardless of how many fsnotify signals fired.

---

## Quick Start

### Prerequisites

- Go 1.21+

### Build

```bash
go build -o cho .
```

### Configure

Copy and edit `cho.toml`:

```toml
# This is the name of your application that calatrava uses to identify a cho client. This is 
# included in the initialHandshake headers. Without it, calatrava will reject the Handshake. 
name = "persona-mp3-cho"

# Cho will tail this file and send the logs to calatrava. This can be the path to 
# any file your application logs to. If you have no logs to point to, run the
# `log_gen.sh` script
logSource = "./logs/structured_logs.txt"

# Address calatrava is running on
ingestorAddr = "http://localhost:9090"

# Interval at which cho will send logs to calatrava. If a logThreshold has been 
# set, logs will only be sent to calatrava when the threshold provided has been met.
interval = "5s"

# Minimum number of logs before cho sends to calatrava. If set to 0, every log 
# created within the interval is sent to calatrava.
logThreshold = 5
```

### Run

```bash
./cho
```

### Generate Test Logs

If you have no application logs to point to:

```bash
# structured JSON logs (recommended)
go run tools/log.go

# simple logs for quick dev iteration
./tools/log_gen.sh

# live feedback loop while developing
watch -n 1 ./tools/log_gen.sh
```

---

## Configuration

| Field | Type | Description |
|---|---|---|
| `name` | string | Identifies this collector to calatrava. Required — handshake is rejected without it |
| `logSource` | string | Path to the log file to tail |
| `ingestorAddr` | string | HTTP address of the calatrava ingestor |
| `interval` | duration | Maximum time between flushes e.g. `"500ms"`, `"5s"` |
| `logThreshold` | int | Minimum buffered logs before flushing. Set to `0` to send every log collected within the interval. Whichever fires first — threshold or interval — triggers a flush |

---

## Log Format

`cho` expects structured JSON logs, one entry per line:

```json
{
  "time": "21-06-2026 12:55:58",
  "level": "ERROR",
  "source": {
    "function": "main.main",
    "file": "/home/james/dev/cho/tools/log.go",
    "line": 85
  },
  "diagnostics": "Generating random logs"
}
```

Parsed into:

```go
type Log struct {
    Timestamp   string
    Level       string
    Diagnostics string
    Source      struct {
        Function string
        File     string
        Line     int
    }
}
```

---

## Design Decisions

**Directory watching over file watching** — `fsnotify` recommends watching directories rather than individual files. On Linux, inotify watches are more reliable at the directory level and correctly handle cases where the file is replaced rather than appended to.

**`[+]` to signal new writes** — borrowed from Vim, which uses `[+]` in the status bar to indicate a buffer has unsaved modifications. Same meaning here: the file has new data.

**ReadAt over Read for offset tracking** — `ReadAt` is a `pread` syscall underneath. It reads from an explicit offset without moving the file pointer, making it safe to call from multiple goroutines and making offset management explicit rather than implicit.

**Threshold + interval dual trigger** — flushing on either condition means high log volume is handled efficiently (threshold triggers before interval, keeps memory bounded) and low log volume still delivers logs promptly (interval triggers before threshold is reached).

**Two-phase handshake with calatrava** — on startup `cho` contacts calatrava before establishing a persistent connection. The handshake negotiates the flush interval so calatrava can tune collector behaviour without redeploying `cho`.

**Single config via TOML** — human-readable, no indentation sensitivity, maps directly to Go structs via BurntSushi/toml.

---

## Known Limitations

**Duplicate log entries** — fsnotify can fire multiple write events for a single logical write depending on how the underlying application flushes. Deduplication before sending to calatrava is not yet implemented.

**At-least-once delivery not guaranteed** — if calatrava is unreachable, buffered logs are dropped. Local disk buffering and retry with backoff are not yet implemented.

**No log rotation handling beyond truncation** — if the log file is replaced entirely (hard rotation), `cho` detects the size decrease and resets the offset to zero. Inode-based rotation detection is not yet implemented.

---

## Status

| Feature | Status |
|---|---|
| TOML config parsing | ✅ |
| Directory watching via fsnotify | ✅ |
| Byte offset tracking | ✅ |
| Structured JSON log parsing | ✅ |
| `logThreshold` + interval batching | ✅ |
| File truncation handling | ✅ |
| Two-phase handshake with calatrava | ✅ |
| HTTP transport to calatrava | 🚧 in progress |
| At-least-once delivery | ⬜ planned |
| Duplicate log deduplication | ⬜ planned |
| Log level filtering | ⬜ planned |
| Multiple log sources | ⬜ planned |
| Log rotation (inode-based) | ⬜ planned |

---

## Project Structure

```
cho/
├── main.go              # entry point, wires config → tailer → watcher
├── cho.go               # Cho struct, readLastLog, offset tracking
├── watcher.go           # fsnotify directory watcher
├── config.go            # TOML config parsing
├── log.go               # Log struct and JSON parsing
├── cho.toml             # default config
└── tools/
    ├── log.go           # structured JSON log generator
    └── log_gen.sh       # simple log generator for dev iteration
```

---

## Later features

- [ ] HTTP transport — send buffered logs to calatrava
- [ ] At-least-once delivery — local buffer when calatrava is unreachable
- [ ] Deduplication — coalesce duplicate log entries before send
- [ ] Inode-based log rotation detection
- [ ] Multiple log sources — tail several files simultaneously
- [ ] Log level filtering — collect only the levels you care about e.g. `ERROR` only
- [ ] Multiple ingestors — route logs from different applications to different calatrava instances
- [ ] stdin support — read from stdout of a process rather than a file

---

## Acknowledgements

- [fsnotify](https://github.com/fsnotify/fsnotify) — cross-platform file system notifications
- [BurntSushi/toml](https://github.com/BurntSushi/toml) — TOML config parsing
- [jkvs](../jkvs/README.md) — the system cho was built to observe
