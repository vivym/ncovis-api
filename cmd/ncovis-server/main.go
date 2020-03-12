package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/getsentry/sentry-go"

	"github.com/vivym/ncovis-api/internal/db"
	"github.com/vivym/ncovis-api/internal/route"
)

// Provisioned by ldflags
// nolint: gochecknoglobals
var (
	version    string
	commitHash string
	buildDate  string
)

func main() {
	v, p := viper.New(), pflag.NewFlagSet(friendlyAppName, pflag.ExitOnError)

	configure(v, p)

	p.String("config", "", "Configuration file")
	p.Bool("version", false, "Show version information")

	_ = p.Parse(os.Args[1:])

	if v, _ := p.GetBool("version"); v {
		fmt.Printf("%s version %s (%s) built on %s\n", friendlyAppName, version, commitHash, buildDate)

		os.Exit(0)
	}

	if c, _ := p.GetString("config"); c != "" {
		v.SetConfigFile(c)
	}

	err := v.ReadInConfig()
	_, configFileNotFound := err.(viper.ConfigFileNotFoundError)
	if !configFileNotFound {
		log.Panic("failed to read configuration", err)
	}

	var config configuration
	err = v.Unmarshal(&config)
	if err != nil {
		log.Panic("failed to unmarshal configuration", err)
	}

	if configFileNotFound {
		log.Println("configuration file not found")
	}

	fmt.Printf("%+v\n", config)

	err = sentry.Init(sentry.ClientOptions{Dsn: config.SentryDsn})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	if err := db.SetupDB(config.DB.URI, config.DB.DBName); err != nil {
		log.Fatalf("mongodb error: %s", err)
	}

	r := route.New(config.GraphiQLToken)
	if err := r.Run("0.0.0.0:" + strconv.Itoa(config.Port)); err != nil {
		log.Fatalf("http error: %s", err)
	}
}
