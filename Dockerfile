FROM golang:latest

RUN mkdir -p /go/src/github.com/LucasDeGaetano/Golang-Challenge
ADD . /go/src/github.com/LucasDeGaetano/Golang-Challenge
WORKDIR /go/src/github.com/LucasDeGaetano/Golang-Challenge
COPY .env . 

RUN go mod download

EXPOSE 8080

CMD ["go", "run", "main.go"]