package main

import (
	"github.com/gin-contrib/cors"
	"github.com/spf13/viper"
)

type Config struct{
	port string
	dbUrl string
	questionPath string
}

func GetConfig() Config {

	viper.SetConfigFile(".env")
    viper.ReadInConfig()

    port := viper.Get("PORT")
    var portString string

    if port == nil {
        portString = ":5000"
    } else {
        portString = port.(string)
    }

    dbUrl := viper.Get("DB_URL")
    var dbUrlString string

    if dbUrl == nil {
        dbUrlString = "quizzie.sqlite.db"
    } else {
        dbUrlString = dbUrl.(string)
    }

	questionPath := viper.Get("QUESTION_PATH")
    var questionPathString string

    if questionPath == nil {
        questionPathString = "../questionPack"
    } else {
        questionPathString = dbUrl.(string)
    }

	return Config{
		port: portString,
		dbUrl: dbUrlString,
		questionPath: questionPathString,
	}
}

func CORSConfig() cors.Config {
    corsConfig := cors.DefaultConfig()
    corsConfig.AllowOrigins = []string{"http://localhost:3000"}
    corsConfig.AllowCredentials = true
    corsConfig.AddAllowHeaders("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers", "Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization")
    corsConfig.AddAllowMethods("GET", "POST", "PUT", "DELETE")
    return corsConfig
}