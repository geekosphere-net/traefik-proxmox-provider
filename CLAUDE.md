# Traefik Proxmox Provider

A Traefik plugin that uses a Proxmox VE cluster as a provider for dynamic configuration.

## Project Structure

- `provider/provider.go` — Main plugin logic: polling, service scanning, Traefik dynamic config generation
- `internal/client.go` — Proxmox API HTTP client
- `internal/types.go` — Data types (Service, ParsedConfig, IP, etc.)
- `go.mod` — Go 1.21+ required (for `log/slog`)

## Key Constraints

- **Yaegi runtime**: This plugin runs inside Traefik's Yaegi Go interpreter. Only standard library packages are available — no third-party logging libraries (zerolog, logrus, zap).
- **`log/slog` for logging**: We use Go's built-in `log/slog` (available since Go 1.21, supported by Yaegi v0.16.0+). The default logger is configured in `provider.New()` based on `config.ApiLogging` ("info" or "debug").
- **stderr for log output**: The slog handler writes to `os.Stderr`, which Traefik captures into its log file. The slog level text in the output (level=ERROR, level=INFO, etc.) provides the actual severity.
- **No `unsafe` package**: Traefik disables `unsafe` for plugins. This doesn't affect stdlib packages like `log/slog` (loaded pre-compiled), but blocks third-party libs that import `unsafe`.

## Logging Conventions

- `slog.Debug()` — Verbose/diagnostic messages (API requests/responses, scanning details, traefik config dumps)
- `slog.Info()` — Normal operational messages (connected, created router, skipping service, hostname fallback)
- `slog.Error()` — Actual errors (panic recovery, config failures, network interface errors)
- Log level is controlled by `ApiLogging` config field: "info" (default) or "debug"

## Build & Test

```bash
make test        # Run unit tests
make lint        # Run golangci-lint
make yaegi_test  # Verify Yaegi interpreter compatibility
make vendor      # Update vendored dependencies
```
