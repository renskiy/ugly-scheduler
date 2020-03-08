package main

import (
	"github.com/renskiy/ugly-scheduler/internal/app"
	"gopkg.in/yaml.v2"
	"os"
)

func main() {
	dbConfig := map[string]*migrationsEnvironment{
		"default": {
			Dialect:    "postgres",
			DataSource: app.DBConnectionString(),
			Dir:        "migrations",
			TableName:  "migrations",
			SchemaName: "public",
		},
	}
	encoder := yaml.NewEncoder(os.Stdout)
	if err := encoder.Encode(dbConfig); err != nil {
		_, _ = os.Stderr.WriteString(err.Error())
	}
}

type migrationsEnvironment struct {
	Dialect    string `yaml:"dialect"`
	DataSource string `yaml:"datasource"`
	Dir        string `yaml:"dir"`
	TableName  string `yaml:"table"`
	SchemaName string `yaml:"schema"`
}
