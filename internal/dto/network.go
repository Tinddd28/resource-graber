package dto

type Network struct {
	Adapters map[string]Adapter
}

type Adapter struct {
	IP      string
	Packets map[string]PacketInfo
}

type PacketInfo struct {
	Protocol    string
	Headers     string
	Source      string
	Destination string
	Sent        int
	Recv        int
	Body        []byte
}
