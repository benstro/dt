# dt

Personal developer tools CLI.

## Install

```bash
go install github.com/benstro/dt@latest
```

## Tools

| Command | Description |
|---|---|
| `dt b64 encode <text>` | Base64 encode |
| `dt b64 decode <text>` | Base64 decode |
| `dt jwt decode <token>` | Decode and pretty-print a JWT header and payload |
| `dt uuid` | Generate a UUID v4 |
| `dt json pretty <json>` | Format JSON with indentation |
| `dt json minify <json>` | Strip whitespace from JSON |
| `dt json stringify <json>` | Escape JSON into a string literal |
| `dt hash sha256 <text>` | SHA-256 hash of a string |
| `dt hash sha256 --file <path>` | SHA-256 hash of a file |
| `dt hash sha256 <file1> <file2> ...` | SHA-256 hash multiple files in parallel |

`b64` and `hash sha256` also accept input from stdin: `echo "hello" | dt b64 encode`

