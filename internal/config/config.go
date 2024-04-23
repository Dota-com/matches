package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	Env      string `yaml:"env"`
	Grpc     GRPC
	Secret   string `yaml:"secret"`
	DevDB    string `yaml:"dev_db"`
	ConfigDB DB     `yaml:"storage_path"`
}

type DB struct {
	Host   string `yaml:"POSTGRES_HOST"`
	UserDb string `yaml:"POSTGRES_USER"`
	DbName string `yaml:"POSTGRES_DB"`
	PassDb string `yaml:"POSTGRES_PASSWORD"`
	PortDb string `yaml:"POSTGRES_PORT"`
}

type GRPC struct {
	Port    int           `yaml:"port"`
	TimeTTl time.Duration `yaml:"timeout"`
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

func MustLoad() *Config {
	fetch := fetchConfigPath()

	if fetch == "" {
		panic("Пустой файл конфига")
	}
	if _, err := os.Stat(fetch); os.IsNotExist(err) {
		panic("Не существует файла конфигурации " + fetch + "\n")
	}

	var conf Config
	if err := cleanenv.ReadConfig(fetch, &conf); err != nil {
		panic("Ошибка чтения конфига " + err.Error() + "\n")
	}

	return &conf
}

func SetupLoger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "dev":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
