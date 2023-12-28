package main

import (
	"fmt"
	ll "lemmy_links_bot/service"
	"os"
	"strconv"
	"strings"

	"github.com/adadkins/glaw"
	"github.com/hashicorp/go-retryablehttp"

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

	config := zap.NewDevelopmentConfig()

	// AddCaller option includes line numbers, file names, and function names
	config.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder

	// Create a Zap logger based on the configuration
	logger, _ := config.Build()

	// configure retryable client
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	client, err := glaw.NewLemmyClient(fmt.Sprintf("%s%s", baseURL, apiVersion), apiToken, jwtCookie, retryClient.StandardClient(), logger)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	a, _ := ll.NewApp(client, logger, banListedAccounts, baseURL)
	a.Work()
}
