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
	Port         string
	JWTSecretKey string
	DbHost       string
	DbPort       int
	DbUser       string
	DbPassword   string
	DbName       string
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
		Port:         os.Getenv("PORT"),
		JWTSecretKey: os.Getenv("JWT_SCECRET"),
		DbHost:       os.Getenv("DB_HOST"),
		DbPort:       dbPort,
		DbUser:       os.Getenv("DB_USER"),
		DbPassword:   os.Getenv("DB_PASSWORD"),
		DbName:       os.Getenv("DB_NAME"),
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
