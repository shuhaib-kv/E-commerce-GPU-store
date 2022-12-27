FROM golang:alpine AS builder 

WORKDIR /E-Commerce-Project/

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ecom .

FROM alpine:latest

WORKDIR /root

COPY --from=builder /E-Commerce-Project/ecom .
     

COPY . .

CMD ["./ecom"]