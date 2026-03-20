run:
	go run cmd/api/main.go

migrate-up:
	goose -dir migrations postgres "host=localhost user=user dbname=dbname password=password sslmode=disable" up

migrate-down:
	goose -dir migrations postgres "host=localhost user=user dbname=dbname password=password sslmode=disable" down

# create a new migration file
# migrate-create:
# 	goose -dir migrations create remove_pgcrypto_extension sql
