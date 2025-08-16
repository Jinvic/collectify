set shell := ["powershell.exe", "-c"]

path := "cmd/... internal/..."

# lint codes [golangci-lint required]
lint:
    golangci-lint cache clean | \
    golangci-lint run {{path}}

# format codes [golangci-lint required]
fmt:
    golangci-lint fmt {{path}}

alias b := build
alias bf := build-frontend
alias bb := build-backend

build:
    cd web | pnpm run build
    go build -o ./collectify.exe .

build-frontend:
    cd web | pnpm run build

build-backend:
    go build -o ./collectify.exe .