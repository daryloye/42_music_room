.SILENT:

.PHONY: run generate-prisma

run:
	cd backend && go run main.go

generate-prisma:
	cd backend && go run github.com/steebchen/prisma-client-go db push
	cd backend && go generate ./...
	cd backend && go mod tidy

create-db:
	docker run --name music-room-postgres \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_PASSWORD=postgres \
    -e POSTGRES_DB=music_room \
    -p 5432:5432 \
    -d postgres:16-alpine
