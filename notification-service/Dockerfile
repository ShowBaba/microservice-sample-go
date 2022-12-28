FROM golang:1.19-alpine
WORKDIR /shared
COPY ../shared .
RUN go mod tidy
WORKDIR /notification-service
COPY ./notification-service .
RUN go mod tidy
RUN go build -o main .
CMD ["./main"]
