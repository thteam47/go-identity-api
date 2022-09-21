FROM golang:latest
WORKDIR /app

COPY . .
RUN go mod download

EXPOSE 8089
EXPOSE 8090

CMD [ "go","run","cmd/main.go" ]