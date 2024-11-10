.PHONY:

build:
	go build -o ./rest cmd/app.go

run: build  
	./rest

swag:
	swag init -g app.go --dir cmd,internal/handlers
swagrm:
	rm -rf ./docs

migrate:
	migrate -path ./migrations -database 'postgresql://postgres:12345@localhost:5430/postgres?sslmode=disable' up
migrate-down:
	migrate -path ./migrations -database 'postgresql://postgres:12345@localhost:5430/postgres?sslmode=disable' down

test:
	go test -v ./tests

up:
	sudo docker compose -f docker-compose.yml up 

