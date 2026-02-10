# logfmt

`logfmt` reads structured logs from stdin and formats into readable output

## Installation

```sh
go install github.com/thdxg/logfmt
```

## Usage

```sh
$ go run ./cmd/autoscaler 2>&1 | logfmt
2026-02-10 12:18:55 INFO Autoscaler started
2026-02-10 12:18:56 DEBUG Scaler registered target=map[id:7bfc2889-8fc2-4b31-b309-80fd8984628b name:RAG]
2026-02-10 12:18:56 INFO Manager started interval=map[metric:5s scale:10s]
2026-02-10 12:19:01 DEBUG Metric fetched target=RAG value=0
2026-02-10 12:19:01 INFO Tick metric success=1 total=1
```

## Motivation

I've been using libraries like [`tint`](https://github.com/lmittmann/tint) to format structured logs in my Go projects.
Formatting logs is primarily for better readability during local development, but using a library for this means adding an unnecessary dependency to your project.
Having a local command line tool to format any kind of json logs solves this problem.

## WIP

- Pretty-printing nested json
- Style configuration
- Handle decoding errors
- Package distribution
