# Stage 1: Build Go application
FROM golang:1.20 as builder

WORKDIR /app

COPY . .

RUN go clean --modcache && \
    go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -o /app/main app/main.go

# Stage 2: Build libwebp from source
FROM alpine:latest AS builder2

RUN apk --no-cache add \
    libpng-dev \
    libjpeg-turbo-dev \
    giflib-dev \
    tiff-dev \
    autoconf \
    automake \
    make \
    gcc \
    g++ \
    wget

RUN wget https://storage.googleapis.com/downloads.webmproject.org/releases/webp/libwebp-0.6.0.tar.gz && \
    tar -xvzf libwebp-0.6.0.tar.gz && \
    mv libwebp-0.6.0 libwebp && \
    rm libwebp-0.6.0.tar.gz && \
    cd /libwebp && \
    ./configure && \
    make && \
    make install && \
    rm -rf libwebp

# Stage 3: Final stage with minimal image
FROM alpine:latest

COPY --from=builder2 /usr/local/bin /usr/local/bin
COPY --from=builder2 /usr/local/include /usr/local/include
COPY --from=builder2 /usr/local/lib /usr/local/lib

RUN apk --no-cache add libpng libjpeg-turbo giflib tiff && \
    rm -rf /usr/local/share /usr/local/libexec

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/public ./public

EXPOSE 8000

CMD ["./main"]
