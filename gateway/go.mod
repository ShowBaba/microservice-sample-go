module github.com/microservice-sample-go/gateway

go 1.19

replace github.com/microservice-sample-go/shared => ../shared

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/joho/godotenv v1.4.0
	github.com/microservice-sample-go/shared v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/rabbitmq/amqp091-go v1.5.0 // indirect
)
