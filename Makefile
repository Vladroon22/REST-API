.PHONY:

build:
	go build -o ./rest cmd/main/app.go

run: build
	./rest

dockers:
	sudo docker pull postgres
	sudo docker build . -t server
docker-rm:
	sudo docker stop server
	sudo docker rm server
	sudo docker stop rest_api_db
	sudo docker rm rest_api_db
docker-rmi:
	sudo docker rmi server
	sudo docker rmi postgres
up:
	sudo docker compose -f docker-compose.yml up