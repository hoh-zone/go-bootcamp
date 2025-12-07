FROM golang:1.22 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chatserver ./cmd/chatserver

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=builder /app/chatserver /app/chatserver
EXPOSE 8082
USER nonroot:nonroot
ENTRYPOINT ["/app/chatserver"]
