package dto

import "time"

type Network struct {
	Adapters map[string]Adapter
}

type Adapter struct {
	// Name string `json:"name"` // EXAMPLE: "eth0"
	IP string `json:"ip"`

	// Data of send and receive: protocol, source, destination, size
	Packets          map[string]PacketInfo `json:"packets"` // key = protocol
	LastActivityTime time.Time
	IPv4             string
	PacketInfo       []PacketInfo
	PacketCounts     map[string]int
	TotalCounts      int
	TotalSent        int
	TotalRecv        int
}

type PacketInfo struct {
	Protocol    string
	Headers     string
	Source      string
	Destination string
	Sent        int
	Recv        int
	Body        []byte
	Timestamp   time.Time
	Direction   string // incoming or outgoind
}
