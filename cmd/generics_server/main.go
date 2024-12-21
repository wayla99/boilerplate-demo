package main

import "boilerplate-demo/src/use_case"

type config struct {
	AppName     string `env:"APP_NAME" envDefault:"boilerplate-demo"`
	AppVersion  string `env:"APP_VERSION" envDefault:"v0.0.0"`
	Environment string `env:"ENVIRONMENT" envDefault:"localhost"`
	Port        uint   `env:"PORT" envDefault:"8080"`
	Debuglog    bool   `env:"DEBUG_LOG" envDefault:"true"`
	Services    struct {
		PostgresqlUri string `env:"DATABASE_DEMO_POSTGRESQL_URI" envDefault:"postgres://postgres:postgres@localhost:5432/demo?sslmode=disable"`
	}
}

func main() {
	cfg := initEnvironment()
	initLogger(cfg)

	useCase := use_case.New(initRepositories(cfg))

	initInterfaces(cfg, useCase)
}
