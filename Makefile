postgres:
	docker run --name postgres14 -p 8080:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:14-alpine

dockerstart:
	docker start postgres14

dockerstop:
	docker stop postgres14

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres14 dropdb simple_bank

migrateup:
	migrate -path db/migration -database postgresql://root:password@localhost:8080/simple_bank?sslmode=disable -verbose up

migratedown:
	migrate -path db/migration -database postgresql://root:password@localhost:8080/simple_bank?sslmode=disable -verbose down

migrateup1:
	migrate -path db/migration -database postgresql://root:password@localhost:8080/simple_bank?sslmode=disable -verbose up 1

migratedown1:
	migrate -path db/migration -database postgresql://root:password@localhost:8080/simple_bank?sslmode=disable -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go -package mockdb simplebank/db/sqlc Store

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

evans:
	evans --host localhost --port 6060 -r repl

.PHONY: postgres dockerstart dockerstop createdb dropdb migrateup migratedown sqlc server mock migrateup1 migratedown1 test proto evans