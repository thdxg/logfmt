# logfmt

[![CI](https://github.com/thdxg/logfmt/actions/workflows/ci.yaml/badge.svg)](https://github.com/thdxg/logfmt/actions/workflows/ci.yaml)
[![Release](https://img.shields.io/github/v/release/thdxg/logfmt)](https://github.com/thdxg/logfmt/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/thdxg/logfmt)](https://goreportcard.com/report/github.com/thdxg/logfmt)

`logfmt` reads structured logs from stdin and formats into readable output

## Installation

### Binary (Recommended)

Download the latest binary for your platform (Linux, macOS, Windows) from the [Releases page](https://github.com/thdxg/logfmt/releases).

### Go Install

```sh
go install github.com/thdxg/logfmt@latest
```

## Usage

```sh
$ go run ./cmd/autoscaler 2>&1 | logfmt
2026-02-10 12:18:55 INFO Autoscaler started
2026-02-10 12:18:56 DEBUG Scaler registered target.id=7bfc2889-8fc2-4b31-b309-80fd8984628b target.name=RAG
2026-02-10 12:18:56 INFO Manager started interval.metric=5s interval.scale=10s
2026-02-10 12:19:01 DEBUG Metric fetched target=RAG value=0
2026-02-10 12:19:01 INFO Tick metric success=1 total=1
```

### Configuration

You can configure `logfmt` using command line flags, environment variables, or a config file (`.logfmt.yaml` in home or current directory).

#### Flags
```sh
logfmt --time-format "15:04:05" --level-format short --color=false
```
- `--time-format`: Timestamp format (Go layout). Default: `2006-01-02 15:04:05`
- `--level-format`: Level style (`full`, `short`, `tiny`). Default: `full`
- `--color`: Enable/disable colored output. Default: `true`
- `--hide-attrs`: Show only time, level, and message. Default: `false`

#### Config File (.logfmt.yaml)
```yaml
time-format: "15:04:05"
level-format: "short"
color: false
hide-attrs: false
```

#### Environment Variables (prefix `LOGFMT_`)
- `LOGFMT_TIME_FORMAT`
- `LOGFMT_LEVEL_FORMAT`
- `LOGFMT_COLOR`
- `LOGFMT_HIDE_ATTRS`

## Motivation

I've been using libraries like [`tint`](https://github.com/lmittmann/tint) to format structured logs in my Go projects.
Formatting logs is primarily for better readability during local development, but using a library for this means adding an unnecessary dependency to your project.
Having a customizable local command line tool to format any kind of json logs solves this problem.
