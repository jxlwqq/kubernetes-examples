FROM golang:1.17-alpine AS builder
WORKDIR /workspace
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install github.com/grpc-ecosystem/grpc-health-probe@latest
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY api api
COPY cmd/calculator cmd/calculator
COPY internal internal
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app cmd/calculator/main.go

FROM alpine
WORKDIR /
COPY --from=builder /workspace/app .
COPY --from=builder /go/bin/grpc-health-probe .
EXPOSE 50051
ENTRYPOINT ["/app"]