package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {

	//Environment Config
	Environment string `mapstructure:"ENVIRONMENT"`
	ServerPort  string `mapstructure:"SERVER_PORT"`

	//Database Config
	DBHost                   string `mapstructure:"DB_HOST"`
	DBPort                   string `mapstructure:"DB_PORT"`
	DBUser                   string `mapstructure:"DB_USER"`
	DBPassword               string `mapstructure:"DB_PASSWORD"`
	DBName                   string `mapstructure:"DB_NAME"`
	DBSSLMode                string `mapstructure:"DB_SSLMODE"`
	DBTimeZone               string `mapstructure:"DB_TIMEZONE"`
	DBMaxIdleConns           int    `mapstructure:"DB_MAX_IDLE_CONNS"`
	DBMaxOpenConns           int    `mapstructure:"DB_MAX_OPEN_CONNS"`
	DBConnMaxLifetimeMinutes int    `mapstructure:"DB_CONN_MAX_LIFETIME_MINUTES"`

	//JWT Config
	JWTSecret            string `mapstructure:"JWT_SECRET"`
	JWTExpirationMinutes int    `mapstructure:"JWT_EXPIRATION_MINUTES"`

	//Email Config
	EmailProvider string `mapstructure:"EMAIL_PROVIDER"`
	EmailAPIKey   string `mapstructure:"EMAIL_API_KEY"`
	EmailSender   string `mapstructure:"EMAIL_SENDER"`

	//Whether the environment is Production,and the default is "-"
	IsProduction bool `mapstructure:"-"`
}

func LoadConfig(path string) (config *Config, errr error) {
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "postgres")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("DB_TIMEZONE", "Asia/Shanghai")
	viper.SetDefault("DB_MAX_IDLE_CONNS", 10)
	viper.SetDefault("DB_MAX_OPEN_CONNS", 100)
	viper.SetDefault("DB_CONN_MAX_LIFETIME_MINUTES", 10)

	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Println("Error reading .env file:", err)
		} else {
			log.Println("No .env file found, will use system env")
			return nil, err
		}
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	config.IsProduction = config.Environment == "production"

	return config, nil
}
