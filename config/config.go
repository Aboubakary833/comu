package config

import (
	"strings"

	"github.com/mazen160/go-random"
	"github.com/spf13/viper"
)

type Config struct {
	AppName      string `mapstructure:"APP_NAME"`
	AppEnv       string `mapstructure:"APP_ENV"`
	AppKey       string `mapstructure:"APP_KEY"`
	AppAddr      string `mapstructure:"APP_ADDR"`
	DBDriver     string `mapstructure:"DB_DRIVER"`
	DBSource     string `mapstructure:"DB_SOURCE"`
	MailHost     string `mapstructure:"MAIL_HOST"`
	MailPort     int    `mapstructure:"MAIL_PORT"`
	MailFrom     string `mapstructure:"MAIL_FROM"`
	MailUserName string `mapstructure:"MAIL_USERNAME"`
	MailPassword string `mapstructure:"MAIL_PASSWORD"`
}

func NewConfig() (*Config, error) {
	var config Config

	setEnvDefaultVariables()

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	err := viper.Unmarshal(&config)

	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(config.AppAddr, ":") {
		config.AppAddr = ":" + config.AppAddr
	}

	return &config, nil
}

func setEnvDefaultVariables() {
	appKey, _ := random.String(64)

	viper.SetDefault("APP_NAME", "Comu")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_ADDR", ":4000")
	viper.SetDefault("APP_KEY", appKey)
	viper.SetDefault("DB_DRIVER", "mysql")
	viper.SetDefault("DB_SOURCE", "root:secret@/comu_db?parseTime=true")
	viper.SetDefault("MAIL_HOST", "localhost")
	viper.SetDefault("MAIL_PORT", "465")
	viper.SetDefault("MAIL_FROM", "norepy@comu.com")
	viper.SetDefault("MAIL_USERNAME", "")
	viper.SetDefault("MAIL_PASSWORD", "")
}
