#build stage
FROM golang:1.18-alpine as builder
WORKDIR /app
COPY . /app
RUN go build -o loggerApp ./main.go
# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/loggerApp .
# COPY wait-for-it.sh .
# RUN apk update && apk add bash
# RUN chmod +x ./wait-for-it.sh
EXPOSE 8002 8002
CMD [ "/app/loggerApp" ]
