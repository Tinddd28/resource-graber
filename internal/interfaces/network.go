package interfaces

import (
	"context"
	"resource-graber/internal/domains/dto"
)

type Network interface {
	Usage(ctx context.Context) dto.Network
}
