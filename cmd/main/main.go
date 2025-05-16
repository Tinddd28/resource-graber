package main

import (
	"log/slog"
	"resource-graber/internal/agents/network"
	"resource-graber/internal/config"
	service "resource-graber/internal/services"
)

func main() {
	config := config.NewConfig()
	// #TODO: write own logger
	logger := slog.Logger{}

	network := network.NewNetwork(&logger)
	service := service.NewService(config, network, nil, nil, nil, &logger)
	service.Run()

}
