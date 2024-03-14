package main

import "github.com/spf13/viper"

// Events Microservice configuration will be read once from .env
type Config struct {
	GinMode         string `mapstructure:"GIN_MODE"`
	Port            string `mapstructure:"SERVICE_PORT"`
	RedisServer     string `mapstructure:"REDIS_SERVER"`
	UseCache        bool   `mapstructure:"USE_CACHE"`
	LogMode         string `mapstructure:"LOG_MODE"`
	ElasticUrl      string `mapstructure:"ELASTICSEARCH_URL"`
	ElasticUser     string `mapstructure:"ELASTIC_USER"`
	ElasticPassword string `mapstructure:"ELASTIC_PASSWORD"`
}

func AppConfig(configFile string) (config Config, err error) {
	viper.SetConfigFile(configFile)
	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	viper.Unmarshal(&config)
	return
}
