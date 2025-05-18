package service

import (
	"context"
	"log/slog"
	"resource-graber/internal/config"
	intrf "resource-graber/internal/interfaces"
)

type Service struct {
	logger  *slog.Logger
	cfg     *config.Config
	network intrf.Network
	client  intrf.ClientAPI
	app     intrf.Application
	screen  intrf.Screen
}

func NewService(cfg *config.Config, network intrf.Network, client intrf.ClientAPI, app intrf.Application, screen intrf.Screen, logger *slog.Logger) *Service {
	return &Service{
		cfg:     cfg,
		network: network,
		client:  client,
		app:     app,
		screen:  screen,
		logger:  logger,
	}
}

func (s *Service) Run(ctx context.Context) {
	s.logger.Info("Service started")
	go s.NetworkUsage(ctx)

}
