FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod tidy

COPY . .

RUN go build -o auth_service .

CMD ["./auth_service"]
