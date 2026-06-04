# dt — Developer Tools CLI

## Glossary

**dt**
The CLI binary. A personal collection of developer tools, installed via `go install github.com/benstro/dt@latest`. Entry point for all tools.

**Tool**
A top-level subcommand of `dt` (e.g. `b64`, `jwt`, `uuid`, `json`). Each tool has its own package under `internal/`.

**Operation**
A sub-subcommand within a tool when the tool has multiple modes (e.g. `encode`/`decode` under `b64`, `pretty`/`minify`/`stringify` under `json`). Single-purpose tools (e.g. `uuid`) have no operation — they run directly.

**b64**
Tool for base64 encoding and decoding. Operations: `encode`, `decode`. Accepts input as an argument or from stdin.

**jwt**
Tool for inspecting JSON Web Tokens. Operation: `decode` — base64-decodes the header and payload and pretty-prints them. No signature verification.

**uuid**
Tool that generates a single random UUID v4.

**json**
Tool for JSON manipulation. Operations: `pretty` (format with indentation), `minify` (strip whitespace), `stringify` (escape a JSON value into a JSON string literal).

**hash**
Tool for hashing input text. Planned — not yet implemented.
