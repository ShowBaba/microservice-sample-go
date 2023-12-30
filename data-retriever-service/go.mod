module github.com/showbaba/microservice-sample-go/data-retriever-service

go 1.19

require (
	github.com/gorilla/mux v1.8.0
	github.com/graphql-go/graphql v0.8.0
	github.com/joho/godotenv v1.4.0
	github.com/showbaba/microservice-sample-go/shared v0.0.0-00010101000000-000000000000
	gorm.io/driver/postgres v1.4.6
	gorm.io/gorm v1.24.3
)

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.2.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/rabbitmq/amqp091-go v1.5.0 // indirect
	golang.org/x/crypto v0.4.0 // indirect
	golang.org/x/text v0.5.0 // indirect
)

replace github.com/showbaba/microservice-sample-go/shared => ../shared
