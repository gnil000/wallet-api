FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .
ADD go.sum .

COPY . .
COPY ./vendor/. ./vendor/.

RUN go build -o main ./cmd/main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/main /main

CMD ["/main"]