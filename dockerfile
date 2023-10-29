FROM alpine:latest

WORKDIR /app
COPY batch-challenge-linux-amd64 /app/

CMD ["/app/batch-challenge-linux-amd64"]
