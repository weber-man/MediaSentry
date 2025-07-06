# Build stage
FROM alpine:latest AS builder

RUN apk add --no-cache gcc go
ENV CGO_ENABLED=1

COPY src/ /app/
WORKDIR /app

RUN go build -v -a -ldflags='-s -w' -o mediasentry .

# Runtime stage
FROM alpine:latest AS runtime

RUN apk add --no-cache ffmpeg
WORKDIR /app

# Copy only the built binary from builder stage
COPY --from=builder /app/mediasentry ./mediasentry
COPY entrypoint.sh ./entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]