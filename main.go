package main

import (
	"datapreparation/app"
	"datapreparation/config"
	"datapreparation/pkg/cryptohelper"
	"datapreparation/pkg/logger"
	"datapreparation/profiles/repository/postgres"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func main() {

	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	runApp()
	viper.OnConfigChange(func(e fsnotify.Event) {
		runApp()
	})
}

func runApp() {

	psDB, err := postgres.NewSQLX(
		viper.GetString(config.PSHost),
		viper.GetString(config.PSPort),
		viper.GetString(config.PSUsername),
		viper.GetString(config.PSDB),
		viper.GetString(config.PSPass),
	)
	if err != nil {
		panic(err)
	}
	decrypt := cryptohelper.NewAESDecryptor([]byte(viper.GetString(config.AesKey)))
	logger, err := logger.InitZap("./app.log")
	if err != nil {
		panic(err)
	}

	app.NewApp(psDB, decrypt).
		WithLogger(logger).
		Run(viper.GetString(config.AppPort))
}
