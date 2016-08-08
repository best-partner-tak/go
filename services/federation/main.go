package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stellar/go/internal/http"
	"github.com/stellar/go/services/federation/internal"
)

var app *federation.App
var rootCmd *cobra.Command

func main() {
	rootCmd.Execute()
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	rootCmd = &cobra.Command{
		Use:   "federation",
		Short: "stellar federation server",
		Long: `stellar federation server
=========================

Make sure config.toml file is in the working folder.
Required config values:
  - domain
  - database.type
  - database.url
  - queries.federation
  - queries.reverse-federation`,
		Run: run,
	}
}

func run(cmd *cobra.Command, args []string) {
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error reading config file: ", err)
	}

	if viper.GetString("database.type") == "" ||
		viper.GetString("database.url") == "" ||
		viper.GetString("domain") == "" ||
		viper.GetString("queries.federation") == "" ||
		viper.GetString("queries.reverse-federation") == "" {
		rootCmd.Help()
		os.Exit(1)
	}

	var config federation.Config
	err = viper.Unmarshal(&config)

	app, err = federation.NewApp(config)

	if err != nil {
		log.Fatal(err)
	}
	http.Run(http.Config{
		ListenAddr: fmt.Sprintf("0.0.0.0:%d", config.Port),
		Handler:    app.Handler(),
	})
}
