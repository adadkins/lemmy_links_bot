package main

import (
	"fmt"
	ll "lemmy_links_bot/service"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/adadkins/glaw"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	apiVersion := os.Getenv("API_VERSION")
	baseURL := os.Getenv("BASE_URL")
	jwtCookie := os.Getenv("JWT_COOKIE")
	apiToken := os.Getenv("API_TOKEN")
	banListedAccounts := []int{}

	for _, str := range strings.Split(os.Getenv("BANLISTED_ACCOUNT_IDS"), ",") {
		intValue, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Error converting '%s' to an integer: %v\n", str, err)
			return
		}
		banListedAccounts = append(banListedAccounts, intValue)
	}
	logger, err := setupLogger()
	if err != nil {
		panic("Error setting up logger: " + err.Error())
	}

	client, err := glaw.NewLemmyClient(fmt.Sprintf("%s%s", baseURL, apiVersion), apiToken, jwtCookie, &http.Client{}, logger)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	a, _ := ll.NewApp(client, logger, banListedAccounts, baseURL)
	a.Work()
}

func getLogLevel() zapcore.Level {
	envLogLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))

	switch envLogLevel {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func setupLogger() (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()

	// Set the logging level based on the environment variable
	config.Level.SetLevel(getLogLevel())

	// AddCaller option includes line numbers, file names, and function names
	config.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder

	// Create a Zap logger based on the configuration
	return config.Build()
}
