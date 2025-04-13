FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=1  \
    GOARCH="amd64" \
    GOOS=linux

WORKDIR /app

RUN apk add --no-cache gcc musl-dev sqlite-dev
COPY . .
RUN go mod tidy
RUN go build --ldflags "-extldflags -static" -o koenbot .


# Stage 2: Minimal runtime image
FROM alpine:latest

WORKDIR /app

# Install cwebp jika butuh untuk sticker conversion
RUN apk add --no-cache ffmpeg

# Copy binary dari builder stage
COPY --from=builder /app/koenbot .

# Temp & session folder
RUN mkdir /app/store /app/temp

VOLUME ["/app/store", "/app/temp"]

# Jalankan bot
CMD ["./koenbot"]