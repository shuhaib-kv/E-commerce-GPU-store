FROM golang:1.12.0-alpine3.9
RUN mkdir /E-com
ADD . /E-com
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./

RUN go build -o /docker-gs-ping
WORKDIR /E-com
RUN go build -o main .
CMD ["/E-com/main"]
