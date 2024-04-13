FROM golang:latest

WORKDIR /REST

RUN apt-get update && apt-get install -y make 

COPY . /REST/

RUN make 

CMD [ "./rest" ]

EXPOSE 8000
