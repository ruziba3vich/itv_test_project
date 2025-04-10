FROM golang:1.24.1 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/main .

EXPOSE 7777

CMD ["./main"]
