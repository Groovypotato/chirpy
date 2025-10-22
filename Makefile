up:
	goose -dir ./sql/schema postgres "postgres://postgres:postgres@localhost:5432/chirpy" up

down:
	goose -dir ./sql/schema postgres "postgres://postgres:postgres@localhost:5432/chirpy" down

connect:
	psql postgres://postgres:postgres@localhost:5432/chirpy