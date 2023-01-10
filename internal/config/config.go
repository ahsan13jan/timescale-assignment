package config

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Config struct {
	Postgres struct {
		Database string `envconfig:"DB" default:"postgres"`
		Hostname string `envconfig:"HOST" default:"127.0.0.1"`
		Port     int    `envconfig:"PORT" default:"5432"`
		Username string `envconfig:"USERNAME" default:"postgres"`
		Password string `envconfig:"PASSWORD" default:"password"`
	} `envconfig:"POSTGRES"`
	CsvPath              string `envconfig:"CSV_PATH"`
	MaxConcurrentWorkers int    `envconfig:"MAX_CONCURRENT_WORKERS" default:"2"`
}

func (c *Config) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		panic(fmt.Errorf("failed to marshal config: %+v", err))
	}
	b = bytes.ReplaceAll(b, []byte("\"Password\":\""+c.Postgres.Password), []byte("\"Password\":\"*****"))

	return string(b)
}
