package main

import (
	"fmt"
	glaw "lemmy_links_bot/lib"
	ll "lemmy_links_bot/service"
	"net/http"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	apiVersion := os.Getenv("API_VERSION")
	baseURL := os.Getenv("BASE_URL")
	jwtCookie := os.Getenv("JWT_COOKIE")
	apiToken := os.Getenv("API_TOKEN")

	config := zap.NewDevelopmentConfig()

	// AddCaller option includes line numbers, file names, and function names
	config.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder

	// Create a Zap logger based on the configuration
	logger, _ := config.Build()

	client, err := glaw.NewLemmyClient(fmt.Sprintf("%s%s", baseURL, apiVersion), apiToken, jwtCookie, &http.Client{}, logger)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	a, _ := ll.NewApp(client, logger)
	a.Work(client, baseURL)
}
