FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o rinha .

EXPOSE 8081
CMD ["/app/rinha"]

