#build stage
FROM golang:1.18-alpine as builder
WORKDIR /app
COPY . /app
RUN go build -o listenApp ./main.go
# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/listenApp .
CMD [ "/app/listenApp" ]
