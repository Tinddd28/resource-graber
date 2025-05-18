package network

import (
	"context"
	"fmt"
	"log"
	"resource-graber/internal/domains/dto"
	"resource-graber/pkg/utils"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
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
	realNet, err := utils.GetRealNetwork()
	if err != nil {
		n.logger.Error("Error getting real network", "error", err)
		return err
	}

	for _, device := range devices {
		n.logger.Debug(fmt.Sprintf("Adapter %s", device.Name))
		ip, ok := realNet[device.Name]
		if ok {
			n.addAdapter(device.Name, ip)
		} else {
			n.logger.Debug(fmt.Sprintf("Skipping adapter %s; not found IP addr", device.Name))
			continue
		}
	}

	return nil
}

func (n *NetworkAgent) addAdapter(name string, ipv4 string) {
	if _, ok := n.adapters[name]; !ok {
		n.adapters[name] = &Adapter{
			IPv4:             ipv4,
			PacketInfo:       make([]dto.PacketInfo, 0),
			mu:               &sync.Mutex{},
			PacketCounts:     make(map[string]int),
			IsUsed:           true,
			IsCapturing:      true,
			LastActivityTime: time.Now(),
			PauseChan:        make(chan string, 1),
		}
		n.logger.Info(fmt.Sprintf("Adding adapter %s", name))
		// n.adapterChan <- AdapterAction{
		// 	Action: "add",
		// 	Name:   name,
		// }
	} else {
		n.adapters[name].mu.Lock()
		n.adapters[name].LastActivityTime = time.Now()
		n.adapters[name].IsUsed = true
		n.adapters[name].IsCapturing = true
		n.adapters[name].mu.Unlock()
	}
}

// All functions will be called when the program starts
// But GetInterfaces() its a main function of this layer. monitorAdapters() will be called once
// like a goroutine in the NewNetwork. They are check the adapters state.
// If the one of them is disconnected, it will be removed from the map, but the data will be saved
// until the capture function process it.`

// #TODO: realize the logic to capture data from the adapter
func (n *NetworkAgent) Capture(ctx context.Context, deviceName string) error {
	time.Sleep(5 * time.Second)
	handle, err := pcap.OpenLive(deviceName, 65535, true, pcap.BlockForever)
	if err != nil {
		return fmt.Errorf("failed to open device %s: %w", deviceName, err)
	}

	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for {
		select {
		case <-ctx.Done():
			n.logger.Info("Data", "ad", n.adapters)
			n.logger.Info("stopping packet capture")
			return nil
		default:
			adapter, exist := n.adapters[deviceName]
			if !exist {
				n.logger.Info(fmt.Sprintf("adapter %s not found", deviceName))
				return nil
			}

			select {
			case command := <-adapter.PauseChan:
				if command == "pause" {
					n.logger.Info(fmt.Sprintf("pausing capture for adapter %s", deviceName))
					adapter.mu.Lock()
					adapter.IsCapturing = false
					adapter.mu.Unlock()
					for {
						select {
						case resumeCommand := <-adapter.PauseChan:
							if resumeCommand == "resume" {
								n.logger.Info(fmt.Sprintf("resuming capture for adapter %s", deviceName))
								adapter.mu.Lock()
								adapter.IsCapturing = true
								adapter.mu.Unlock()
								break
							}
						case <-ctx.Done():
							n.logger.Info("stopping packet capture")
							return nil
						}
					}
				}
			default:

			}

			select {
			case packet, ok := <-packetSource.Packets():
				if !ok {
					n.logger.Info(fmt.Sprintf("packet source closed for adapter %s", deviceName))
					return nil
				}

				n.processPacket(packet, deviceName)
			case <-ctx.Done():
				n.logger.Info("stopping packet capture")
				return nil
			}

		}
	}
}

func (n *NetworkAgent) processPacket(packet gopacket.Packet, deviceName string) {
	adapter, exists := n.adapters[deviceName]
	if !exists {
		return
	}

	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		return
	}
	ip, _ := ipLayer.(*layers.IPv4)

	var protocol string
	switch {
	case packet.Layer(layers.LayerTypeTCP) != nil:
		protocol = TCPProtocol
	case packet.Layer(layers.LayerTypeUDP) != nil:
		protocol = UDPProtocol
	case packet.Layer(layers.LayerTypeICMPv4) != nil:
		protocol = ICMPProtocol
	default:
		protocol = UnknownValue
	}
	var direction string
	var source, destination string
	if adapter.IPv4 == ip.SrcIP.String() {
		direction = Outgoing
		source = ip.SrcIP.String()
		destination = ip.DstIP.String()
	} else if adapter.IPv4 == ip.DstIP.String() {
		direction = Incoming
		source = ip.DstIP.String()
		destination = ip.SrcIP.String()
	} else {
		return
	}
	if packet == nil {
		return
	}

	domain := extractDNS(packet)
	if domain == "" {
		return
	}
	n.logger.Info(domain)

	packetInfo := dto.PacketInfo{
		Source:      source,
		Destination: destination,
		Direction:   direction,
		Protocol:    protocol,
		Sent:        len(packet.Data()),
		Recv:        len(packet.Data()),
		Body:        packet.Data(),
		Timestamp:   time.Now(),
	}
	// n.logger.Info(packet.String())

	adapter.mu.Lock()
	defer adapter.mu.Unlock()
	adapter.PacketInfo = append(adapter.PacketInfo, packetInfo)
	adapter.TotalCounts.Add(1)
	// #FIXME: need grouping by remote IP
	adapter.PacketCounts[domain]++
	if direction == Outgoing {
		adapter.TotalSent.Add(int64(len(packet.Data())))
	} else {
		adapter.TotalRecv.Add(int64(len(packet.Data())))
	}

}

func (n *NetworkAgent) monitorAdapters(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 5)
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
			var wg sync.WaitGroup
			for name, adapter := range n.adapters {
				wg.Add(1)
				go func(name string, adapter *Adapter) {
					defer wg.Done()

					adapter.mu.Lock()
					defer adapter.mu.Unlock()
					if curDevice[name] {
						adapter.LastActivityTime = time.Now()
						adapter.IsUsed = true
						if !adapter.IsCapturing {
							adapter.PauseChan <- "resume"
						}
					} else {
						// if time.Since(adapter.LastActivityTime) > inactiveTimeout {
						adapter.IsUsed = false
						adapter.PauseChan <- "pause"
						// }
					}
				}(name, adapter)
			}

			for _, device := range devices {
				if _, exists := n.adapters[device.Name]; !exists {
					ip := utils.GetIPv4(device)
					if ip != "unknown" {
						n.addAdapter(device.Name, ip)
					}
				}
			}
		}
	}
}

// func extractSNI(packet gopacket.Packet) string {
// 	tlsLayer := packet.Layer(layers.LayerTypeTLS)
// 	if tlsLayer == nil {
// 		return ""
// 	}

// 	tls, ok := tlsLayer.(*layers.TLS)
// 	if !ok || len(tls.Handshake) == 0{
// 		return ""
// 	}

// 	for _, record := range tls.Handshake {
// 		if record.ContentType == layers.TLSHandshake {
// 			cl :=
// 		}
// 	}

// }

func extractDNS(packet gopacket.Packet) string {
	dnsLayer := packet.Layer(layers.LayerTypeDNS)
	if dnsLayer != nil {
		return ""
	}

	dns, ok := dnsLayer.(*layers.DNS)
	if !ok || dns == nil {
		return ""
	}
	log.Println(dns)
	if len(dns.Questions) > 0 {
		return string(dns.Questions[0].Name)
	}
	return ""
}
