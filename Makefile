.PHONY: run
run:
	go run cmd/main.go

.PHONY: migrate-up
migrate-up:
	migrate -source file:$(PWD)/migrations -database "postgres://postgres:mysecretpassword@localhost:5432/mydb?sslmode=disable" up

.PHONY: migrate-down
migrate-down:
	migrate -source file:$(PWD)/migrations -database "postgres://postgres:mysecretpassword@localhost:5432/mydb?sslmode=disable" down