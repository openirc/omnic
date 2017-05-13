package main

import (
    "fmt"

    "github.com/spf13/viper"
)

func read_config() {
    viper.SetConfigName("omnic")
    viper.AddConfigPath("./conf")
    viper.AddConfigPath(".")

    err := viper.ReadInConfig()
    if err != nil {
        panic(fmt.Errorf("Fatal error config file: %s \n", err))
    }

    viper.SetDefault("httpd.bindip", "0.0.0.0")
    viper.SetDefault("httpd.bindport", "8080")
    viper.SetDefault("psql.connectstring", "db=openirc")
}
