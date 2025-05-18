package main

import (
	// "context"
	"context"
	"log/slog"
	"time"

	// "resource-graber/internal/agents/network"
	"resource-graber/internal/agents/network"
	"resource-graber/internal/config"
	service "resource-graber/internal/services"
	"resource-graber/pkg/logger"
	// "time"
)

func main() {
	config := config.NewConfig()
	// #TODO: write own logger

	logger := logger.SetupLogger("local")
	logger.Info("Starting network usage monitoring")
	logger.Info("", slog.Any("cfg: ", config))
	network := network.NewNetwork(logger)
	service := service.NewService(config, network, nil, nil, nil, logger)
	// network.GetAdaptersList()
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Minute*3))
	defer cancel()
	service.Run(ctx)
	// network.GetInterfaces()
	// network.Capture(ctx, "wlp4s0")
	time.Sleep(5 * time.Minute)

	// handle, err := pcap.OpenLive("wlp4s0", 65535, true, pcap.BlockForever)
	// if err != nil {
	// 	panic(err)
	// }
	// defer handle.Close()

	// packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

}
