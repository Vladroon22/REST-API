.PHONY:

build:
	go build -o ./rest cmd/main/app.go

run: build # work only with postgres image
	./rest

test:
	go test -v ./tests

dockers:
	sudo docker pull postgres:16.2
	sudo docker build . -t rest-api-server
docker-rm:
	sudo docker stop server
	sudo docker rm server
	sudo docker stop rest_api_db
	sudo docker rm rest_api_db
docker-rmi:
	sudo docker rmi rest-api-server
	sudo docker rmi postgres:16.2
up:
	sudo docker compose -f docker-compose.yml up

# tests in docker
# 	sudo docker exec -it server sh
#	go -v ./tests