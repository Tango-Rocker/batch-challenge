FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
COPY resources/basic_schema.json .

RUN go mod download
COPY . .
RUN go build -o main .
#EXPOSE 8080

# Command to run the executable
CMD ["./main"]
