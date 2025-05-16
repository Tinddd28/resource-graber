package network

import (
	"context"
	"fmt"
	"resource-graber/internal/dto"
	"sync"
	"time"

	"github.com/google/gopacket/pcap"
)

func (n *NetworkAgent) GetInterfaces() error {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return err
	}

	if len(devices) == 0 {
		return fmt.Errorf(zeroDeviceValue)
	}
	// #TODO: realize the logic to add adapters
	// this function will be called when the program starts
	// and will add all available adapters to the list
	// for _, device := range devices {
	// }

	return nil
}

func (n *NetworkAgent) addAdapter(name string, ip string) {
	if _, ok := n.adapters[name]; !ok {
		n.adapters[name] = &Adapter{
			IP:           ip,
			PacketInfo:   make([]dto.PacketInfo, 0),
			mu:           &sync.Mutex{},
			PacketCounts: make(map[string]int),
			IsUsed:       true,
		}
	}
}

// All functions will be called when the program starts
// But GetInterfaces() its a main function of this layer. monitorAdapters() will be called once
// like a goroutine in the NewNetwork. They are check the adapters state.
// If the one of them is disconnected, it will be removed from the map, but the data will be saved
// until the capture function process it.`

// #TODO: realize the logic to capture data from the adapter
// this function will be called when
func (n *NetworkAgent) Capture(ctx context.Context, deviceName string) error {
	return nil
}

func (n *NetworkAgent) monitorAdapters(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			n.logger.Info("Stopping network adapter monitoring")
			return
		case <-ticker.C:
			devices, err := pcap.FindAllDevs()
			if err != nil {
				n.logger.Error("Error finding devices", "error", err)
				continue
			}

			curAdapters := make(map[string]bool, len(n.adapters))
			for name := range n.adapters {
				curAdapters[name] = true
			}
		}
	}
}
