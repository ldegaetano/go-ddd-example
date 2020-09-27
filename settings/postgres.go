package settings

import "github.com/kelseyhightower/envconfig"

type postgresSettings struct {
	DBName   string `envconfig:"DB_NAME" required:"true"`
	UserName string `envconfig:"DB_USER_NAME" required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
	Host     string `envconfig:"DB_HOST" required:"true"`
	Port     string `envconfig:"DB_PORT" required:"true"`
}

var Postgres postgresSettings

func init() {
	if err := envconfig.Process("", &Postgres); err != nil {
		panic(err.Error())
	}
}
