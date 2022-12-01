FROM golang:1.19.3-alpine3.16 as builder

WORKDIR /app

# Download deps
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build
COPY *.go .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "bin/app"

# Deploy
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /app/bin/app .
USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["./app"]