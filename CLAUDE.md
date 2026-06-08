# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go build ./...          # compile everything
go run .                # run without installing
go install .            # install dt to $GOPATH/bin
go test ./...           # run all tests
go test ./internal/b64/ # run tests for a single package
go mod tidy             # sync go.mod/go.sum after adding/removing dependencies
```

## Architecture

`main.go` delegates entirely to `cmd.Execute()`. All Cobra wiring lives in `cmd/`, all logic lives in `internal/`.

**Adding a new tool** follows this pattern every time:
1. Create `internal/<tool>/<tool>.go` — pure functions, no Cobra, no fmt.Println, no os.Exit
2. Create `cmd/<tool>.go` — import the internal package, wire up Cobra commands in `init()`

`init()` in each `cmd/*.go` file calls `rootCmd.AddCommand(...)` — this is how tools self-register without `main.go` needing to know about them.

**Two-level command structure:** top-level tools (`b64`, `jwt`, `json`) have subcommand operations (`encode`/`decode`, `pretty`/`minify`/`stringify`). Single-purpose tools (`uuid`) run directly with no subcommand.

**Input resolution:** tools accept input as a CLI argument or from stdin (pipe). The `resolveInput` helper in `cmd/input.go` handles this by checking `os.ModeCharDevice` on stdin's stat — use this pattern for any tool that should support piped input.

## Domain

See `CONTEXT.md` for the glossary of terms (Tool, Operation, etc.) used throughout this project.

The `hash` tool (`hash sha256`) hashes a string, a single file (`--file <path>`), or multiple files (positional args, hashed in parallel with goroutines). Output for files is `<hash>  <path>`, matching `sha256sum` convention.
