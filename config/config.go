package config

import "github.com/jinzhu/configor"

// Config is the application configuration placeholder.
type Config struct {
	AppName    string `default:"go-api"`
	AppVersion string `required:"true"`
	AppHost    string `default:"127.0.0.1"`
	AppPort    string `default:"3000"`

	Database struct {
		Client   string `default:"postgres"`
		Host     string `default:"127.0.0.1"`
		Name     string `required:"true"`
		User     string `required:"true"`
		Password string `required:"true"`
		Port     string `default:"5432"`
	}

	RabbitMQ struct {
		Host     string `default:"127.0.0.1"`
		Name     string `required:"true"`
		User     string `required:"true"`
		Password string `required:"true"`
		Port     string `default:"5672"`
	}

	Logger struct {
		Host  string
		Port  string `default:"12201"`
		Level string `default:"INFO"`
	}
}

// Load config
func (c *Config) Load(fileName string) error {
	return configor.Load(c, fileName)
}
