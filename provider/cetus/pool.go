package cetus

import (
	"math/big"

	"github.com/omnibtc/go-sui-swap/types"
)

type cetusPool struct {
	poolAddress  string
	coinAAddress string
	coinBAddress string
	config       *CetusPoolConfig
}

func (p *cetusPool) Address() string {
	return p.poolAddress
}

func (p *cetusPool) PoolType() types.SuiPoolType {
	return types.PoolTypeCetus
}

func (p *cetusPool) CoinA() string {
	return p.coinAAddress
}

func (p *cetusPool) CoinB() string {
	return p.coinBAddress
}

func (p *cetusPool) GetAmountOut(_ types.Pool, inAsset types.Asset, outAsset types.Asset, inAmount *big.Int) (*big.Int, error) {
	return p.config.Quoter.GetAmountOut(p, inAsset, outAsset, inAmount)
}
