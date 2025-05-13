FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY ../src/go.mod ../src/go.sum ./
RUN go mod download
COPY ../src .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bikeboxes-collector .

# final stage
FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/bikeboxes-collector .
RUN adduser -D -g '' appuser
USER appuser

# running the binary
ENTRYPOINT ["./bikeboxes-collector"]