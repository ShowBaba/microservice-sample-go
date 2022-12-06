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
	MailUsername string
	MailPassword string
}

func GetConfig() Config {
	if _config == nil {
		_config = ConfigFactory()
	}
	return *_config
}

func defaultConfig() *Config {
	return &Config{
		MailUsername: os.Getenv("MAIL_USERNAME"),
		MailPassword: os.Getenv("MAIL_PASSWORD"),
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
