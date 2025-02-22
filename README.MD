# My Go Web project

## Useful command

### Go

- get dependencies:
  `go get github.com/joho/godotenv`
- copy dependencies into vendor to avoid download from remote:
  `go mod vendor`
- clean up:
  `go mod tidy`
- build & run:
  `go build && ./my-first-go`
- install cli tool
  `go install github.com/pressly/goose/v3/cmd/goose@latest`

### SQL (with `sqlc` & `goose`)

- up migration (perform "Up" action in migration sql file):
  `goose postgres [DB_URL] up`
- down:
  `goose postgres [DB_URL] down`
- generate go-typed - like prisma typescript generate
  `sqlc generate`
