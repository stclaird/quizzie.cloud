package main

import (
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
        portString = ":3000"
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
