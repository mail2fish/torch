package utilities

import (
	"fmt"

	"github.com/go-playground/log"
	"github.com/spf13/viper"
)

var env string

// GetConfigString 根据环境变量的服务器环境，获取配置文件里的配置
func GetConfigString(key string) string {

	return viper.GetString(fmt.Sprintf("%s.%s", env, key))
}

// GetConfigString 根据环境变量的服务器环境，获取配置文件里的配置列表
func GetConfigStringSlice(key string) []string {
	return viper.GetStringSlice(fmt.Sprintf("%s.%s", env, key))
}

func InitEnv() {
	env = viper.GetString("ENV")
	log.Info("Environment is ", env)
}
