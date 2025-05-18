package network

import (
	"context"
	"resource-graber/internal/domains/dto"
	"sync"
)

func (n *NetworkAgent) Usage(ctx context.Context) dto.Network {
	network := dto.Network{
		Adapters: make(map[string]dto.Adapter),
	}
	go n.monitorAdapters(ctx)
	for name, _ := range n.adapters {
		go n.Capture(ctx, name)
	}

	for {
		select {
		case <-ctx.Done():
			var wg sync.WaitGroup
			for name, adapter := range n.adapters {
				wg.Add(1)
				go func(name string, adapter *Adapter) {
					defer wg.Done()
					adapter.mu.Lock()

					network.Adapters[name] = dto.Adapter{
						IPv4:         adapter.IPv4,
						PacketCounts: adapter.PacketCounts,
						TotalCounts:  int(adapter.TotalCounts.Load()),
						TotalSent:    int(adapter.TotalSent.Load()),
						TotalRecv:    int(adapter.TotalRecv.Load()),
						PacketInfo:   adapter.PacketInfo,
					}

					// #TODO: need clear unused device

					adapter.PacketInfo = make([]dto.PacketInfo, 0)
					adapter.PacketCounts = make(map[string]int)
					adapter.TotalCounts.Store(0)
					adapter.TotalSent.Store(0)
					adapter.TotalRecv.Store(0)

					adapter.mu.Unlock()
				}(name, adapter)
			}

			wg.Wait()
			return network
		}
	}
}

func (n *NetworkAgent) GetAdaptersList() {
	for name, adapter := range n.adapters {
		n.logger.Info(name, "ip", adapter.IPv4)
	}
}
