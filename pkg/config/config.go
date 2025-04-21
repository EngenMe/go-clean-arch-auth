package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	JWTSecret  string `mapstructure:"JWT_SECRET"`
	ServerPort string `mapstructure:"SERVER_PORT"`
	Env        string `mapstructure:"APP_ENV"`
}

var (
	config Config
	once   sync.Once
)

func LoadConfig() (Config, error) {
	var err error

	once.Do(
		func() {
			env := os.Getenv("APP_ENV")
			if env == "" {
				env = "development"
			}

			v := viper.New()
			v.SetConfigType("yaml")

			v.SetConfigName("config")
			v.AddConfigPath("./configs")
			if err = v.ReadInConfig(); err != nil {
				log.Printf(
					"Warning: Could not read base config file: %v\n",
					err,
				)
			}

			v.SetConfigName(fmt.Sprintf("config.%s", env))
			v.AddConfigPath("./configs")
			if err = v.MergeInConfig(); err != nil {
				log.Printf(
					"Warning: Could not read environment config file: %v\n",
					err,
				)
			}

			v.AutomaticEnv()

			v.SetDefault("DB_HOST", "localhost")
			v.SetDefault("DB_PORT", "5432")
			v.SetDefault("DB_USER", "respositories")
			v.SetDefault("DB_PASSWORD", "respositories")
			v.SetDefault("DB_NAME", "auth_db")
			v.SetDefault("JWT_SECRET", "your-secret-key")
			v.SetDefault("SERVER_PORT", "8080")
			v.SetDefault("APP_ENV", env)

			// Unmarshal into config struct
			if err = v.Unmarshal(&config); err != nil {
				log.Printf("Error unmarshaling config: %v\n", err)
			}

			config.Env = env
		},
	)

	return config, err
}

func GetConfig() Config {
	_, err := LoadConfig()
	if err != nil {
		log.Printf("Error loading configs: %v\n", err)
	}
	return config
}
