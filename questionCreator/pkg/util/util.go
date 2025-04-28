package util

import (
	"github.com/gin-contrib/cors"
	"github.com/spf13/viper"
)

type Config struct{
	Port string
	Questionsdir string
	Geminimodel string
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func CORSConfig() cors.Config {
    corsConfig := cors.DefaultConfig()
    corsConfig.AllowOrigins = []string{"http://localhost:5001"}
    corsConfig.AllowCredentials = true
    corsConfig.AddAllowHeaders("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers", "Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization")
    corsConfig.AddAllowMethods("GET", "POST", "PUT", "DELETE")
    return corsConfig
}
