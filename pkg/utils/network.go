package utils

import (
	"net"

	"github.com/google/gopacket/pcap"
)

func GetIPv4(device pcap.Interface) string {
	for _, addr := range device.Addresses {
		if ipv4 := addr.IP.To4(); ipv4 != nil {
			return ipv4.String()
		}
	}

	return "unknown"
}

func GetRealNetwork() (map[string]string, error) {
	res := make(map[string]string)
	intrfs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, intrf := range intrfs {
		if intrf.Flags&net.FlagLoopback != 0 || intrf.Flags&net.FlagUp == 0 {
			continue
		}

		ips, err := intrf.Addrs()
		if err != nil {
			return nil, err
		}

		for _, ip := range ips {
			if ipnet, ok := ip.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				res[intrf.Name] = ipnet.IP.String()
			}
		}
	}

	return res, nil
}
