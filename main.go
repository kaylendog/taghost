package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
)

var log *zap.SugaredLogger

func init() {
	viper.SetConfigName("taghost")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	viper.SetDefault("server.port", 8080)
	viper.SetDefault("assets.path", "./assets")

	viper.SetDefault("git.repository_url", "https://github.com/username/repository")
	viper.SetDefault("git.username", "username")
	viper.SetDefault("git.access_token", "token")

	viper.SetDefault("database.username", "username")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.dbname", "dbname")
	viper.SetDefault("database.host", "host")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = viper.SafeWriteConfig()
			if err != nil {
				panic("failed to write config")
			}
		} else {
			panic("failed to read config")
		}
	}
}

func main() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()

	if err != nil {
		panic(err)
	}

	log = logger.Sugar()
	defer logger.Sync()

	db, err := gorm.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.username"),
		viper.GetString("database.dbname"),
		viper.GetString("database.password"),
	))

	if err != nil {
		panic(err)
	}

	defer db.Close()
	log.Info("Created DB connection pool")

	CheckAndClone()

	router := NewRouter()

	log.Info(fmt.Sprintf("Listening on :%s", viper.GetString("server.port")))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", viper.GetString("server.port")), router))
}
