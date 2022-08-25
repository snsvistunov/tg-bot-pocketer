package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken  string
	ConsumerKey    string
	AuthServerURL  string
	TelegramBotURL string `mapstructure:"bot_url"`
	DBPath         string `mapstructure:"db_file"`

	Messages Messages

	Commands Commands
}

type Messages struct {
	Errors
	Responses
}

type Errors struct {
	Default      string `mapstructure:"default"`
	InvalidURL   string `mapstructure:"invalid_url"`
	Unauthorized string `mapstructure:"unauthorized"`
	UnableToSave string `mapstructure:"unable_to_save"`
}

type Responses struct {
	Start             string `mapstructure:"start"`
	AlreadyAuthorized string `mapstructure:"already_authorized"`
	SavedSuccessfully string `mapstructure:"saved_successfully"`
	UnknownCommand    string `mapstructure:"unknown_command"`
}

type Commands struct {
	Start string `mapstructure:"start"`
}

func Load() (*Config, error) {
	var cfg Config

	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.errors", &cfg.Messages.Errors); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.responses", &cfg.Messages.Responses); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("commands", &cfg.Commands); err != nil {
		return nil, err
	}

	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseEnv(cfg *Config) error {

	if err := viper.BindEnv("TOKEN"); err != nil {
		return err
	}

	if err := viper.BindEnv("CONSUMER_KEY"); err != nil {
		return err
	}

	if err := viper.BindEnv("AUTH_SERVER_URL"); err != nil {
		return err
	}

	cfg.TelegramToken = viper.GetString("TOKEN")
	cfg.ConsumerKey = viper.GetString("CONSUMER_KEY")
	cfg.AuthServerURL = viper.GetString("AUTH_SERVER_URL")

	return nil
}
