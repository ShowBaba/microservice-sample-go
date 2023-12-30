package app

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var (
	_config       *Config
	ConfigFactory = defaultConfig
)

type Config struct {
	Port              string
	MongoURI          string
	RabbitmqServerURL string
}

func GetConfig() Config {
	if _config == nil {
		_config = ConfigFactory()
	}

	return *_config
}

func defaultConfig() *Config {
	return &Config{
		Port:              os.Getenv("PORT"),
		MongoURI:          os.Getenv("MONGO_URI"),
		RabbitmqServerURL: os.Getenv("RABBITMQ_SERVER_URL"),
	}
}

func init() {
	var (
		dir, _   = os.Getwd()
		basepath = filepath.Join(dir, ".env")
	)
	log.Println("env basepath: ", basepath)
	if err := godotenv.Load(basepath); err != nil {
		log.Print("No .env file found")
		panic(err)
	}
}
