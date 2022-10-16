postgres:
	docker run --name postgres14 -p 8080:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:14-alpine

startpostgres:
	docker start postgres14

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres14 dropdb simple_bank

migrateup:
	migrate -path db/migration -database postgresql://root:password@localhost:8080/simple_bank?sslmode=disable -verbose up

migratedown:
	migrate -path db/migration -database postgresql://root:password@localhost:8080/simple_bank?sslmode=disable -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go -package mockdb simplebank/db/sqlc Store

.PHONY: postgres startpostgres createdb dropdb migrateup migratedown sqlc server mock