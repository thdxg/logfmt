# Implementation Plan: Support slog TextHandler (logfmt)

## Goal
Enable `logfmt` to parse logs formatted by Go's `slog.TextHandler` (logfmt style: `key=value`) in addition to JSON.

## Approach
Since the user requested **no new dependencies**, we will implement a custom parser for the logfmt format.

## Steps

1.  **Create `pkg/parser` Package**
    *   Create `pkg/parser/logfmt.go`.
    *   Implement `ParseLogfmt(line []byte) (map[string]any, error)`.
    *   **Parsing Logic**:
        *   Iterate through the input string.
        *   Extract keys (text before `=`).
        *   Extract values:
            *   If value starts with `"`, scan until the closing quote (handling escaped quotes). Use `strconv.Unquote` to decode.
            *   Otherwise, scan until the next space.
        *   Store pairs in a `map[string]any`.

2.  **Add Tests**
    *   Create `pkg/parser/logfmt_test.go`.
    *   Test cases:
        *   Simple: `time=now level=INFO msg=hello`
        *   Quoted: `msg="hello world" key="value with spaces"`
        *   Escaped: `msg="quoted \"inside\""`
        *   Mixed: `a=1 b="2" c=3`
        *   Invalid inputs (graceful handling).

3.  **Integrate into `cmd/root.go`**
    *   Modify the main loop in `Run`.
    *   Current logic: `json.Unmarshal` -> success? -> format.
    *   New logic:
        *   Try `json.Unmarshal`.
        *   If it fails, try `parser.ParseLogfmt`.
        *   If that succeeds (returns valid map), use it.
        *   If both fail, print original line.

4.  **Verification**
    *   Run tests.
    *   Manual verification with `slog` output examples.

## Verification
- Run `go test ./...`
- Pipe `slog` output to `logfmt`.
