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

- cd to schema folder: `cd sql/schema`
  - up migration (perform "Up" action in migration sql file):
    `goose postgres [DB_URL] up`
  - down:
    `goose postgres [DB_URL] down`
- without cd to schema folder: setting `GOOSE_*` in .env
  `goose up`
- generate go-typed - like prisma typescript generate
  `sqlc generate`

## Note:

- this project follows [youtube tutorial](https://www.youtube.com/watch?v=un6ZyFkqFKo) but have some changes due to some dependencies out of date
- use `pgx/v5` lastest
- use `go-chi/chi/v5`
- use native `net/http` instead of `go-chi`: [branch `native-net/http`](../../tree/native-net/http) - manual: https://www.youtube.com/watch?v=H7tbjKFSg58
