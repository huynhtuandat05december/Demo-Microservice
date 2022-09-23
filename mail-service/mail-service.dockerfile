#build stage
FROM golang:1.18-alpine as builder
WORKDIR /app
COPY . /app
RUN go build -o mailApp ./main.go
# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/mailApp .
RUN mkdir /templates
COPY templates templates
EXPOSE 8003 8003
CMD [ "/app/mailApp" ]
