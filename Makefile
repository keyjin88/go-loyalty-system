.PHONY: build

build:
	go build -v ./cmd/gophermart

test:
	go test ./...

migrate_up:
	migrate -path database/migrations/ -database "postgresql://pgadmin:postgres@localhost:5432/loyaltydb?sslmode=disable" -verbose up
migration_down:
	migrate -path database/migrations/ -database "postgresql://pgadmin:postgres@localhost:5432/loyaltydb?sslmode=disable" -verbose down
migration_fix:
	migrate -path database/migration/ -database "postgresql://pgadmin:postgres@localhost:5432/loyaltydb?sslmode=disable" force VERSION

.DEFAULT_GOAL:= build