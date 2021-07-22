FROM golang:1.16-alpine AS builder
WORKDIR /go/src/hello-gin
ADD . /go/src/hello-gin
RUN go get -d -v ./...
RUN go build -o app .

FROM alpine:latest
COPY --from=builder /go/src/hello-gin/app .
EXPOSE 8080
CMD ["./app"]