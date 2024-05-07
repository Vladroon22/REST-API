.PHONY:

build:
	go build -o ./rest cmd/main/app.go

run: build 
# 	in configfile -> conf.toml change host = localhost and port = 5430 
#	sudo docker pull postgres:16.2 
#	sudo docker run --name=rest_api_db -e POSTGRES_PASSWORD=12345 -p 5430:5432 -d postgres:16.2
	./rest

test:
	go test -v ./tests

up:
	sudo docker compose -f docker-compose.yml up 
#	in configfile -> conf.toml change host = postgres-db

tests-in-docker:
	sudo docker exec -it server sh
#	make test

docker-rm:
	sudo docker stop server
	sudo docker rm server
	sudo docker stop rest_api_db
	sudo docker rm rest_api_db
docker-rmi:
	sudo docker rmi rest-api-server
	sudo docker rmi postgres:16.2