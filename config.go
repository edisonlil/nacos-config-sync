package main

import (
	"github.com/spf13/viper"
	"os"
	"regexp"
	"strings"
)

type Environment struct {
	Config *viper.Viper
}

var env *Environment

func init() {

	config := viper.GetViper()
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

	env = &Environment{
		Config: config,
	}

	env.bindEnvs()
}

func (env *Environment) GetString(key string) string {
	return env.resolvePlaceholders(env.Config.GetString(key))
}

func (env *Environment) bindEnvs() {

	env.Config.BindEnv("NACOS_USERNAME")
	env.Config.BindEnv("NACOS_PASSWORD")
}

func (env *Environment) resolvePlaceholders(val string) string {

	//Match %\w*% 匹配%{任意字符或数字或下划线组成的单词}%
	if ok, _ := regexp.Match("\\$\\{\\w*}", []byte(val)); !ok {
		return val
	}
	return env.Config.GetString(strings.ReplaceAll(strings.ReplaceAll(val, "${", ""), "}", ""))
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

func (env *Environment) GetNacosConfig() NacosConfig {
	return NacosConfig{
		Addr:          env.GetString("sync-config.nacos.addr"),
		Url:           env.GetString("sync-config.nacos.url"),
		Group:         env.GetString("sync-config.nacos.group"),
		Namespace:     env.GetString("sync-config.nacos.namespace"),
		FileExtension: env.GetString("sync-config.nacos.file-extension"),
		LoginUrl:      env.GetString("sync-config.nacos.loginUrl"),
		ConfigUrl:     env.GetString("sync-config.nacos.configUrl"),
		Username:      env.GetString("sync-config.nacos.username"),
		Password:      env.GetString("sync-config.nacos.password"),
	}
}
