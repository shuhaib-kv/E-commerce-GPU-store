FROM golang:alpine AS builder 

WORKDIR /E-COMMERCE-PROJECT/

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ecom_api .

FROM alpine:latest

WORKDIR /root

COPY --from=builder /E-COMMERCE-PROJECT/ecom_api .
     

COPY . .

CMD ["./ecom_api"]