#build stage
FROM golang:1.18-alpine as builder
WORKDIR /app
COPY . /app
RUN go build -o brokerApp ./cmd/api
# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/brokerApp .
EXPOSE 80 80
CMD [ "/app/brokerApp" ]
