postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine
createdb :
	docker exec -it postgres14 createdb --username=root --owner=root simple_bank
dropdb :
	docker exec -it postgres14 dropdb simple_bank
migrateup:
	migrate -path db/migration -database="postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database="postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
dockerps :
	docker ps -a
test :
	go test -v -cover ./...
sqlc :
	sqlc generate
.PHONY:postgres createdb dropdb dockerps migrateup migratedown sqlc test