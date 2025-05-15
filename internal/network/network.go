package network

import (
	"log/slog"
	intrf "resource-graber/internal/interfaces"
)

type Network struct {
	logger  *slog.Logger
	Path    string
	Network intrf.Network
	Client  intrf.ClientAPI
}

func NewNetwork(path string, network intrf.Network, client intrf.ClientAPI, logger *slog.Logger) *Network {
	return &Network{
		Path:    path,
		Network: network,
		Client:  client,
		logger:  logger,
	}
}
