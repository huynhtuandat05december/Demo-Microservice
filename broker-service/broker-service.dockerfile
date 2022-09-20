#build stage
FROM golang:1.18-alpine as builder
WORKDIR /app
COPY . /app
RUN go build -o brokerApp ./main.go
# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/brokerApp .
EXPOSE 8000 8000
CMD [ "/app/brokerApp" ]
