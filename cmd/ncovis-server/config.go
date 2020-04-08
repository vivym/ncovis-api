package main

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type configuration struct {
	SentryDsn string

	Port int

	DB struct {
		DBName string
		URI    string
	}

	GraphiQLToken string
	AdminToken    string
}

// configure configures some defaults in the Viper instance.
func configure(v *viper.Viper, p *pflag.FlagSet) {
	// Viper settings
	v.AddConfigPath(".")
	v.AddConfigPath(fmt.Sprintf("$%s_CONFIG_DIR/", strings.ToUpper(envPrefix)))

	// Environment variable settings
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AllowEmptyEnv(true)
	v.AutomaticEnv()

	// Application constants
	v.Set("appName", appName)

	v.SetDefault("sentryDsn", "https://da71503fc8de4f41909fe1a3539ae53f@sentry.io/3831823")
	v.SetDefault("port", 12711)
	v.SetDefault("db.dbname", "ncovis")
	v.SetDefault("db.uri", "mongodb://localhost:27017/")
	v.SetDefault("GraphiQLToken", "")
	v.SetDefault("AdminToken", "")
}
