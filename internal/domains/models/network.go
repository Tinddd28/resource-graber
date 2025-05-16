package models

type Network struct {
	Adapter map[string]Adapter `json:"adapters"` //key = name of adapter
}

type Adapter struct {
	// Name string `json:"name"` // EXAMPLE: "eth0"
	IP string `json:"ip"`

	// Data of send and receive: protocol, source, destination, size
	Packets   map[string]PacketInfo `json:"packets"` // key = protocol
	TotalSent int                   `json:"total_sent"`
	TotalRecv int                   `json:"total_recv"`
}

type PacketInfo struct {
	Headers     string `json:"headers"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Sent        int    `json:"send_bytes"`
	Recv        int    `json:"receive_bytes"`
	Body        []byte `json:"body"`
}
