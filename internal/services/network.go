package service

import (
	"context"
	"time"
)

func (s *Service) NetworkUsage(glCtx context.Context) {
	s.logger.Info("Starting network usage monitoring")
	ctx := context.Background()
	ticker := time.NewTicker(s.cfg.Network.Timeout)
	for {
		select {
		case <-glCtx.Done():
			s.logger.Info("Stopping network usage monitoring")
			return
		case <-ticker.C:
			ctxTimeout, cancel := context.WithTimeout(ctx, 20*time.Second)
			defer cancel()
			network := s.network.Usage(ctxTimeout)
			for name, adapter := range network.Adapters {
				s.logger.Info(name, "ipv4", adapter.IPv4)
				for k, v := range adapter.PacketCounts {
					s.logger.Info("PacketCounts", k, v)
				}
				s.logger.Info("TotalSent", "bytes", adapter.TotalSent)
				s.logger.Info("TotalRecv", "bytes", adapter.TotalRecv)
			}
		}
	}
}
