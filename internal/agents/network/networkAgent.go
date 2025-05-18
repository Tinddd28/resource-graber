package network

import (
	"log/slog"
	"resource-graber/internal/domains/dto"
	"sync"
	"sync/atomic"
	"time"
)

type Adapter struct {
	LastActivityTime time.Time
	IPv4             string
	PacketInfo       []dto.PacketInfo
	mu               *sync.Mutex
	PacketCounts     map[string]int // domain name: count
	TotalCounts      atomic.Int64
	TotalSent        atomic.Int64
	TotalRecv        atomic.Int64
	IsUsed           bool        // flag to check if the adapter still connected
	IsCapturing      bool        // flag to check if the adapter is capturing packets
	PauseChan        chan string // channel to pause capture values:pause, resume
}

type AdapterAction struct {
	Name   string
	Action string // add, remove
}

// #TODOL need think about how to use this field, what we can add or remove from this struct
type NetworkAgent struct {
	logger   *slog.Logger
	adapters map[string]*Adapter
	doneChan chan struct{}
	// adapterChan chan AdapterAction
}

func NewNetwork(logger *slog.Logger) *NetworkAgent {

	na := &NetworkAgent{
		logger:   logger,
		adapters: make(map[string]*Adapter),
		// adapterChan: make(chan AdapterAction),
	}

	na.GetInterfaces()

	return na
}
