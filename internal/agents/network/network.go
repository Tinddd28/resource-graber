package network

import (
	"context"
	"fmt"
	"resource-graber/internal/domains/dto"
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
	for _, device := range devices {
		n.addAdapter(device.Name)
	}

	return nil
}

func (n *NetworkAgent) addAdapter(name string) {
	if _, ok := n.adapters[name]; !ok {
		n.adapters[name] = &Adapter{
			PacketInfo:   make([]dto.PacketInfo, 0),
			mu:           &sync.Mutex{},
			PacketCounts: make(map[string]int),
			IsUsed:       true,
		}

		n.adapterChan <- AdapterAction{
			Action: "add",
			Name:   name,
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

func (n *NetworkAgent) monitorAdapters(ctx context.Context, inactiveTimeout time.Duration) {
	ticker := time.NewTicker(time.Second)
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

			curDevice := make(map[string]bool, len(devices))
			for _, device := range devices {
				curDevice[device.Name] = true
			}

			for name, adapter := range n.adapters {
				if curDevice[name] {
					adapter.mu.Lock()
					adapter.LastActivityTime = time.Now()
					adapter.IsUsed = true
					adapter.mu.Unlock()
				} else {
					adapter.mu.Lock()
					if time.Since(adapter.LastActivityTime) > inactiveTimeout {
						adapter.IsUsed = false
					}
					adapter.mu.Unlock()
				}
			}

			for _, device := range devices {
				if _, exists := n.adapters[device.Name]; !exists {
					n.addAdapter(device.Name)
				}
			}
		}
	}
}
