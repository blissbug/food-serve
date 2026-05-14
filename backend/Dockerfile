FROM ubuntu:latest
LABEL authors="Blisss"

#base image
FROM golang:1.25.5

#working directory
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

# copy the folder
COPY . .

# run stuff inside the container
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping ./cmd/api/main.go

CMD ["/docker-gs-ping"]