FROM golang:latest

RUN apt-get update && apt-get install -y make 

COPY . /REST
COPY go.mod /REST
COPY go.sum /REST

WORKDIR /REST

RUN make 
CMD [ "./rest" ]