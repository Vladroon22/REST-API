FROM golang:1.22.1

WORKDIR /REST

RUN apt-get update && apt-get install -y make 

COPY . /REST/

RUN go mod download
RUN make 

CMD [ "./rest" ]

EXPOSE 8000
