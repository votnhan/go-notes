FROM golang:1.16-alpine

RUN apk update && apk add gcc g++ sqlite 

RUN mkdir /init_db && mkdir /install

COPY db/*  /init_db

# log db
RUN chmod 777 /init_db/init_db.sh && /init_db/init_db.sh

COPY go.mod /install/go.mod
COPY go.sum /install/go.sum
RUN cd /install && go mod download

WORKDIR /usr/local/app

COPY . .

RUN go build main.go 
RUN go build -o consumer cmd/log_consumer.go 
