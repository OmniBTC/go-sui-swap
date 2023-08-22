package deepbook

import (
	"math/big"

	"github.com/omnibtc/go-sui-swap/types"
)

type DeepBookPool struct {
	poolAddress string
	baseAsset   string
	quoteAsset  string
	quoter      types.Quoter
}

func NewDeepBookPool(poolAddress, baseAsset, quoteAsset string, quoter types.Quoter) types.Pool {
	return &DeepBookPool{
		poolAddress: poolAddress,
		baseAsset:   baseAsset,
		quoteAsset:  quoteAsset,
		quoter:      quoter,
	}
}

func (p *DeepBookPool) Address() string {
	return p.poolAddress
}

func (p *DeepBookPool) CoinA() string {
	return p.baseAsset
}

func (p *DeepBookPool) CoinB() string {
	return p.quoteAsset
}

func (p *DeepBookPool) PoolType() types.SuiPoolType {
	return types.PoolTypeDeepBook
}

func (p *DeepBookPool) GetAmountOut(_ types.Pool, inAsset types.Asset, outAsset types.Asset, inAmount *big.Int) (*big.Int, error) {
	return p.quoter.GetAmountOut(p, inAsset, outAsset, inAmount)
}
