package postgres

import "fmt"

type Config struct {
	Host     string            `mapstructure:"host" envconfig:"HOST" default:"127.0.0.1"`
	Port     int               `mapstructure:"port" envconfig:"PORT" default:"5432"`
	Database string            `mapstructure:"database" envconfig:"DATABASE" default:"postgres"`
	User     string            `mapstructure:"user" envconfig:"USER" default:"postgres"`
	Password string            `mapstructure:"password" envconfig:"PASSWORD" default:"postgres"`
	Params   map[string]string `mapstructure:"params" envconfig:"PARAMS" default:"sslmode:require"`
}

// DSN assembles the connection string for the given Config.
// e.g. postgres://user:password@host:port/database?param1=value1
func (c Config) DSN() string {
	const format = "postgres://%s:%s@%s:%d"
	str := fmt.Sprintf(format, c.User, c.Password, c.Host, c.Port)

	if c.Database != "" {
		str += fmt.Sprintf("/%s", c.Database)
	}

	var params string
	for key, value := range c.Params {
		params += fmt.Sprintf("%s=%s", key, value)
	}

	if params != "" {
		str += fmt.Sprintf("?%s", params)
	}

	return str
}
