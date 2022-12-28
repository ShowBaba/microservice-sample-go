FROM golang:1.19-alpine
WORKDIR /shared
COPY ../shared .
RUN go mod tidy
WORKDIR /auth-service
COPY ./auth-service .
RUN go mod tidy
RUN go build -o main .
EXPOSE 3000
CMD ["./main"]
