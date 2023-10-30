FROM alpine:latest

WORKDIR /app
COPY exec_app /app/
COPY resources/basic_schema.json /app/

CMD ["/app/exec_app"]
