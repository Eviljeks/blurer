.PHONY: server-run
server-run:
	DATABASE_URL="postgres://blurer:blurer@localhost:5432/blurer" go run cmd/server/main.go

.PHONY: db-migrate
db-migrate:
	migrate -source file://migrations/ -database "postgres://blurer:blurer@localhost:5432/blurer?sslmode=disable" up