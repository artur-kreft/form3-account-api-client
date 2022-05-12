FROM golang:1.18 AS build

# Add required packages
RUN apt-get update
RUN apt-get install curl
RUN apt-get install git
RUN apt-get install bash

WORKDIR /app

COPY . .
RUN ls -a
RUN go mod tidy
RUN go mod download

#RUN go test ./...
