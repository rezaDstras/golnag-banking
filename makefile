postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine
createdb :
	docker exec -it postgres14 createdb --username=root --owner=root simple_bank
dropdb :
	docker exec -it postgres14 dropdb simple_bank
migrateup:
	migrate -path db/migration -database="postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
lastmigrationup:
	migrate -path db/migration -database="postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database="postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
lastmigratedown:
	migrate -path db/migration -database="postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
dockerps :
	docker ps -a
test :
	go test -v -cover ./...
sqlc :
	sqlc generate
server :
	go run main.go
.PHONY:postgres createdb dropdb dockerps migrateup lastmigrationup migratedown lastmigratedown sqlc test server