FROM alpine:latest

WORKDIR /app
COPY batch-challenge-linux-amd64 /app/
COPY resources/schema.json /app/

# Set environment variable for schema path
ENV SCHEMA_PATH=/app/schema.json

CMD ["/app/batch-challenge-linux-amd64"]
