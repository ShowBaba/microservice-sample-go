FROM golang:1.19-alpine
WORKDIR /shared
COPY ../shared .
RUN go mod tidy
WORKDIR /gateway-service
COPY ./gateway-service .
RUN go mod tidy
RUN go build -o main .
EXPOSE 3002
CMD ["./main"]
