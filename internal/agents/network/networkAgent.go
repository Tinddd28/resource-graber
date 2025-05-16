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
	PacketInfo       []dto.PacketInfo
	mu               *sync.Mutex
	PacketCounts     map[string]int
	TotalCounts      atomic.Int64
	TotalSent        atomic.Int64
	TotalRecv        atomic.Int64
	IsUsed           bool // flag to check if the adapter still connected
}

type AdapterAction struct {
	Name   string
	IP     string
	Action string // add, remove
}

// #TODOL need think about how to use this field, what we can add or remove from this struct
type NetworkAgent struct {
	logger      *slog.Logger
	adapters    map[string]*Adapter
	doneChan    chan struct{}
	adapterChan chan AdapterAction
}

func NewNetwork(logger *slog.Logger) *NetworkAgent {
	return &NetworkAgent{
		logger:   logger,
		adapters: make(map[string]*Adapter),
	}
}
