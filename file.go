package main

import (
	"log"

	"github.com/spf13/viper"
)

func viperEnvVariable() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Fatalf("No se encontro archivo config.yaml %s", err)
		} else {
			// Config file was found but another error was produced
			log.Fatalf("Hay un error con el archivo config.yaml %s", err)
		}
	}
}
