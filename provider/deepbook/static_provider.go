package deepbook

import (
	"context"

	"github.com/omnibtc/go-sui-swap/types"
)

type staticProvider struct {
	pools []types.Pool
}

func NewStaticProvider(pools []types.Pool) types.Provider {
	return &staticProvider{
		pools: pools,
	}
}

func (p *staticProvider) Pools(ctx context.Context) ([]types.Pool, error) {
	return p.pools, nil
}
