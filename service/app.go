package lemmylinks_service

import (
	"github.com/adadkins/glaw"
	"go.uber.org/zap"
)

type App struct {
	banListedAccounts []int
	done              chan struct{}
	lemmyClient       glaw.Client
	logger            *zap.Logger
}

func NewApp(lc glaw.Client, logger *zap.Logger, banListedAccounts []int) (*App, error) {
	// TODO: how do i properly shut this down?
	done := make(chan struct{})

	return &App{
		banListedAccounts: banListedAccounts,
		done:              done,
		lemmyClient:       lc,
		logger:            logger,
	}, nil
}
