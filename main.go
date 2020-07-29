package main

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
)

func init() {
	viper.SetConfigName("taghost")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	viper.SetDefault("server.port", 8080)
	viper.SetDefault("assets.path", "./assets")

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
	cfg := zap.Config{
		Level:    zap.NewAtomicLevelAt(zap.DebugLevel),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := cfg.Build()
	defer logger.Sync()
	log := logger.Sugar()

	router := NewRouter()

	log.Info(fmt.Sprintf("Listening on :%s", viper.GetString("server.port")))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", viper.GetString("server.port")), router))
}
