.PHONY:

build:
	go build -o ./rest cmd/main/app.go

run: build  
#	sudo docker pull postgres:16.2 
#	sudo docker run --name=postgres -e POSTGRES_PASSWORD=12345 -p 5430:5432 -d postgres:16.2
	./rest

migrate:
	migrate -path ./migrations -database 'postgresql://postgres:12345@localhost:5430/postgres?sslmode=disable' up
migrate-down:
	migrate -path ./migrations -database 'postgresql://postgres:12345@localhost:5430/postgres?sslmode=disable' down

test:
	go test -v ./tests

up:
	sudo docker compose -f docker-compose.yml up 

tests-in-docker:
	sudo docker exec -it server sh
#	make test

docker-rm:
	sudo docker rm -f server
	sudo docker rm -f postgres
docker-rmi:
	sudo docker rmi rest-api-server
	sudo docker rmi postgres:16.2