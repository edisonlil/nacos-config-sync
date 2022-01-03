package main

import (
	"github.com/spf13/viper"
	"os"
	"strings"
)

var config *viper.Viper

func InitConfig() {

	config = viper.GetViper()
	config.SetConfigName("sync-config")
	config.SetConfigType("yml")

	var configDir string

	if len(os.Args) == 3 && os.Args[1] == "-c" {
		configDir = os.Args[2]
	} else {
		configDir, _ = os.Getwd()
	}

	config.AddConfigPath(configDir)
	config.ReadInConfig()
}

func GetConfigEnv() *viper.Viper {
	envConfig := &viper.Viper{}
	envConfig.SetConfigType("yml")
	envConfig.ReadConfig(strings.NewReader(config.GetString("sync-config.env")))
	return envConfig
}

type NacosConfig struct {
	Addr          string
	Url           string
	Group         string
	Namespace     string
	FileExtension string
	LoginUrl      string
	ConfigUrl     string
	Username      string
	Password      string
}

func GetNacosConfig() NacosConfig {
	return NacosConfig{
		Addr:          config.GetString("sync-config.nacos.addr"),
		Url:           config.GetString("sync-config.nacos.url"),
		Group:         config.GetString("sync-config.nacos.group"),
		Namespace:     config.GetString("sync-config.nacos.namespace"),
		FileExtension: config.GetString("sync-config.nacos.file-extension"),
		LoginUrl:      config.GetString("sync-config.nacos.loginUrl"),
		ConfigUrl:     config.GetString("sync-config.nacos.configUrl"),
		Username:      config.GetString("sync-config.nacos.username"),
		Password:      config.GetString("sync-config.nacos.password"),
	}
}
