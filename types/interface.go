package types

import (
	"context"
	"math/big"
)

type SuiPoolType int

const (
	PoolTypeCetus SuiPoolType = iota
	PoolTypeDeepBook
)

type Quoter interface {
	GetAmountOut(pool Pool, inAsset Asset, outAsset Asset, inAmount *big.Int) (*big.Int, error)
}

type Asset interface {
	Address() string
}

type Pool interface {
	Quoter
	Address() string
	PoolType() SuiPoolType
	CoinA() string
	CoinB() string
}

type Provider interface {
	Pools(ctx context.Context) ([]Pool, error)
}

type Trade interface {
	Path() []Asset
	Pools() []Pool
	AmountOuts() []*big.Int
	A2B() []bool
}
