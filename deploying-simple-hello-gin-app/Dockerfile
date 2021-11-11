# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM golang:1.17-alpine AS builder
ARG TARGETARCH
ARG TARGETOS
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -o app main.go

FROM gcr.io/distroless/static:nonroot-$TARGETARCH
WORKDIR /
COPY --from=builder /workspace/app .
EXPOSE 8080
ENTRYPOINT ["/app"]