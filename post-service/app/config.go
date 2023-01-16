package app

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	_config       *Config
	ConfigFactory = defaultConfig
)

type Config struct {
	Port                 string
	DbHost               string
	DbPort               int
	DbUser               string
	DbPassword           string
	DbName               string
	BlogServiceSecretKey string
	RabbitmqServerURL    string
	JWTSecretKey         string
}

func GetConfig() Config {
	if _config == nil {
		_config = ConfigFactory()
	}

	return *_config
}

func defaultConfig() *Config {
	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	return &Config{
		Port:                 os.Getenv("PORT"),
		DbHost:               os.Getenv("DB_HOST"),
		DbPort:               dbPort,
		DbUser:               os.Getenv("DB_USER"),
		DbPassword:           os.Getenv("DB_PASSWORD"),
		DbName:               os.Getenv("DB_NAME"),
		BlogServiceSecretKey: os.Getenv("BLOG_SERVICE_SECRET_KEY"),
		RabbitmqServerURL:    os.Getenv("RABBITMQ_SERVER_URL"),
		JWTSecretKey:         os.Getenv("JWT_SCECRET"),
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
