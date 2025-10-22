db_string := "postgres://postgres:postgres@localhost:5432/chirpy"
db_test_string := "postgres://postgres:postgres@localhost:5432/chirpy_test"

connect:
	psql $(db_string)

up:
	goose -dir ./sql/schema postgres $(db_string) up

down:
	goose -dir ./sql/schema postgres $(db_string) down






connect_test:
	psql $(db_test_string)

up_test:
	goose -dir ./sql/schema postgres $(db_test_string) up

down_test:
	goose -dir ./sql/schema postgres $(db_test_string) down



clean:
	go clean -cache -modcache -testcache
	go mod tidy
	go build ./...
