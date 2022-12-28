FROM golang:1.19-alpine
WORKDIR /shared
COPY ../shared .
RUN go mod tidy
WORKDIR /user-service
COPY ./user-service .
RUN go mod tidy
RUN go build -o main .
EXPOSE 3003
CMD ["./main"]
