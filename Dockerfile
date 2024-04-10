#syntax=docker/dockerfile:1

# stage 1
FROM golang:1.20 as builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go clean --modcache
RUN go mod tidy 
RUN CGO_ENABLED=0 GOOS=linux go build app/main.go

# stage 2
FROM alpine:3
WORKDIR /root/
COPY --from=builder /app/main .
# RUN mkdir -p ./ssl
# RUN touch ./ssl/certificate.crt
# RUN touch ./ssl/private.key
EXPOSE 8000
CMD ["./main"]
