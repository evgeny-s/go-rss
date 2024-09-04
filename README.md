# The concurrent RSS aggregator

### Quick start
```bash
cp .env.example .env
go build && ./go-rss
```


### Run DB migrations:
```bash
cd sql/schema
goose postgres postgres://go-rss-user:go-rss-password@localhost:5439/go-rss-db up
```

### Generate go code using SQL:
```bash
cd <root-folder>
sqlc generate
```

### Run the DB using docker compose:
```bash
export .env
docker-compose up
```