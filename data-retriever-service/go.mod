module github.com/showbaba/microservice-sample-go/data-retriever-service

go 1.19

require (
	github.com/graphql-go/graphql v0.8.0
	github.com/joho/godotenv v1.4.0
	github.com/showbaba/microservice-sample-go/shared v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/rabbitmq/amqp091-go v1.5.0 // indirect
)

replace github.com/showbaba/microservice-sample-go/shared => ../shared
