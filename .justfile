set shell := ["powershell.exe", "-c"]
set dotenv-load

path := "cmd/... internal/..."

# lint codes [golangci-lint required]
lint:
    golangci-lint cache clean | \
    golangci-lint run {{path}}

# format codes [golangci-lint required]
fmt:
    golangci-lint fmt {{path}}

# show version [golangci-lint required]
version:
    golangci-lint --version