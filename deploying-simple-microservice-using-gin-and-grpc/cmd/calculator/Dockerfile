FROM golang:1.17-alpine AS builder

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY api api
COPY cmd/calculator cmd/calculator

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app cmd/calculator/main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/app .
EXPOSE 50051
ENTRYPOINT ["/app"]