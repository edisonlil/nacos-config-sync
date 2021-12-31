package main

import (
	"github.com/spf13/viper"
	"os"
	"strings"
)
var config *viper.Viper

func InitConfig(){

	config = viper.GetViper()
	config.SetConfigName("sync-config")
	config.SetConfigType("yml")

	var configDir string

	if len(os.Args) == 3 && os.Args[1] == "-c" {
		configDir = os.Args[2]
	}else {
		configDir,_ = os.Getwd()
	}

	config.AddConfigPath(configDir)
	config.ReadInConfig()
}

func GetConfigEnv() *viper.Viper{
	envConfig := &viper.Viper{}
	envConfig.SetConfigType("yml")
	envConfig.ReadConfig(strings.NewReader(config.GetString("sync-config.env")))
	return envConfig
}