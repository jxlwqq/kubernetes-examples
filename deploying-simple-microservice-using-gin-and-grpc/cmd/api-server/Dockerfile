FROM golang:1.17-alpine AS builder
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY api api
COPY cmd/api-server cmd/api-server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app cmd/api-server/main.go

FROM alpine
WORKDIR /
COPY --from=builder /workspace/app .
EXPOSE 8080
ENTRYPOINT ["/app"]