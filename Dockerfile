FROM golang:alpine as builder
WORKDIR /E-Commerce-Project

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN  go build -o ecom main.go

FROM alpine:latest

WORKDIR /root
COPY --from=builder /E-Commerce-Project/ecom .
COPY . .
EXPOSE 8080
CMD ["./ecom"]

