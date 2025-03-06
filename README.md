## About

A basic web subnet lookup tool

### Configuration
The following configuration methods are accepted, in order of highest to lowest priority:
- Command-line flags
- Environment variables

### Environment variables
Almost all options configurable via flags can also be configured via environment variables. 

The associated environment variable is the prefix `TRIVIA_` plus the flag name, with the following changes:
- Leading hyphens removed
- Converted to upper-case
- All internal hyphens converted to underscores

For example:
- `--bind 127.0.0.1` becomes `SUBNETTER_BIND=127.0.0.1`
- `--profile` becomes `SUBNETTER_PROFILE=true`

## Usage output
```
Serves a tool for learning IP subnetting.

Usage:
  subnetter [flags]

Flags:
  -b, --bind string       address to bind to (default "0.0.0.0")
      --exit-on-error     shut down webserver on error, instead of just printing the error
  -h, --help              help for subnetter
  -p, --port uint16       port to listen on (default 8080)
      --profile           register net/http/pprof handlers
      --tls-cert string   path to TLS certificate
      --tls-key string    path to TLS keyfile
  -v, --verbose           log requests to stdout
  -V, --version           display version and exit
  ```