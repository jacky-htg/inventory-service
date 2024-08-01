FROM golang:alpine

WORKDIR /app
COPY inventory-service /app/
COPY ./cmd/cmd-inventory-service /app/

CMD ["/app/inventory-service"]
