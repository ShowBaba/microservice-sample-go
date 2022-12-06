module github.com/microservice-sample-go/notification-service

go 1.19

require github.com/joho/godotenv v1.4.0

require (
	github.com/microservice-sample-go/shared v0.0.0-00010101000000-000000000000
	github.com/rabbitmq/amqp091-go v1.5.0
)

require github.com/lib/pq v1.10.7 // indirect

replace github.com/microservice-sample-go/shared => ../shared
