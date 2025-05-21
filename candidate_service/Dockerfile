FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o candidate_service .

EXPOSE 5003

CMD ["./candidate_service"]
