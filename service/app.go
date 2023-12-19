package lemmylinks_service

import (
	glaw "lemmy_links_bot/lib"

	"go.uber.org/zap"
)

type App struct {
	done        chan struct{}
	lemmyClient glaw.Client
	logger      *zap.Logger
}

func NewApp(lc glaw.Client, logger *zap.Logger) (*App, error) {
	// TODO: how do i properly shut this down?
	done := make(chan struct{})

	return &App{
		done:        done,
		lemmyClient: lc,
		logger:      logger,
	}, nil
}
