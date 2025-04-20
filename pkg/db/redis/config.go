package redis

import "fmt"

type Config struct {
	Hosts     string `envconfig:"HOSTS" default:"127.0.0.1" required:"true"`
	Port      int    `envconfig:"PORT" default:"6379" required:"true"`
	Username  string `envconfig:"USERNAME"`
	Password  string `envconfig:"PASSWORD"`
	Database  int    `envconfig:"DATABASE" default:"0"`
	IsCluster bool   `envconfig:"DATABASE" default:"false"`
}

func (c Config) DSN() string {
	dsn := "redis://"

	if c.Username != "" && c.Password != "" {
		dsn += fmt.Sprintf("%s:%s@", c.Username, c.Password)
	} else if c.Username != "" {
		dsn += fmt.Sprintf("%s@", c.Username)
	} else if c.Password != "" {
		dsn += fmt.Sprintf(":%s@", c.Password)
	}

	dsn += fmt.Sprintf("%s:%d", c.Hosts, c.Port)

	if c.Database == 0 {
		return dsn
	}

	return dsn + fmt.Sprintf("/%d", c.Database)
}
