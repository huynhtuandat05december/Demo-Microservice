#build stage
FROM golang:1.18-alpine as builder
WORKDIR /app
COPY . /app
RUN go build -o authenticationApp ./main.go
# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/authenticationApp .
COPY wait-for-it.sh .
RUN apk update && apk add bash
RUN chmod +x ./wait-for-it.sh
EXPOSE 8001 8001
CMD [ "/app/authenticationApp" ]
