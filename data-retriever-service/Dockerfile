FROM golang:1.19-alpine
WORKDIR /shared
COPY ../shared .
RUN go mod tidy
WORKDIR /data-retriever-service
COPY ./data-retriever-service .
RUN go mod tidy
RUN go build -o main .
EXPOSE 3005
CMD ["./main"]
