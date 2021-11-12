package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Configurations struct {
	Server                 ServerConfigurations
	ExchangeRateService    ExchangeRateServiceConfigurations
	ExchangeRateController ExchangeRateControllerConfigurations
}

type ExchangeRateServiceConfigurations struct {
	Urls []string
}

type ExchangeRateControllerConfigurations struct {
	CacheDurationMinutes int
}

type ServerConfigurations struct {
	Port string
}

func NewConfigurations() (config *Configurations) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return config
}
