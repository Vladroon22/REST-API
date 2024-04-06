.PHONY:

build:
	go build -o ./rest cmd/main/app.go

run: build
	./rest

# docker build

