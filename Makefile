.PHONY:

build:
	go build -o ./rest cmd/main/app.go

run: build
	./rest

docker:
	docker build . -t server
	docker run --name=server -p 8080:8000 -p 5430:5432 -d server

docker-rm:
	docker stop server
	docker rm server

docker-rmi:
	docker rmi server

up:
	docker-compose up 