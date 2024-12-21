FROM golang:1.21 AS builder
WORKDIR /app
ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /go/bin/app cmd/generics_server/*.go

FROM alpine:3.15

WORKDIR /app
RUN apk --no-cache add ca-certificates

ARG IMAGE_TAG
ENV APP_VERSION=$IMAGE_TAG

COPY --from=builder /go/bin/app /app/app

CMD ["/app/app"]
