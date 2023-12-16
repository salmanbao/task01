migrate:
	migrate create -ext sql -dir db/migrations -seq init_schema

migrate-up:
	migrate -path db/migrations -database "postgres://salmansaleem:root@localhost:5432/user_data?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migrations -database "postgres://salmansaleem:root@localhost:5432/user_data?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: sqlc