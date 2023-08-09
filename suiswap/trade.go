package suiswap

import (
	"math/big"
	"sync"

	"github.com/omnibtc/go-sui-swap/types"
)

type innerTrade struct {
	pools     []types.Pool
	path      []types.Asset
	isA2B     []bool
	amountOut *big.Int

	wg *sync.WaitGroup
}

func newTrade(pools []types.Pool, path []string, isA2B []bool, amountIn *big.Int, wg *sync.WaitGroup) types.Trade {
	assets := make([]types.Asset, len(path))
	for _, p := range path {
		assets = append(assets, types.NewAsset(p))
	}

	innerTrade := &innerTrade{
		pools: pools,
		isA2B: isA2B,
		path:  assets,
		wg:    wg,
	}
	wg.Add(1)
	go innerTrade.asyncTradeAmountIn(amountIn, wg)
	return innerTrade
}

func (t *innerTrade) asyncTradeAmountIn(amountIn *big.Int, wg *sync.WaitGroup) {
	defer wg.Done()
	var err error
	for i, p := range t.pools {
		amountIn, err = p.GetAmountOut(t.path[i], t.path[i+1], amountIn)
		if err != nil {
			return
		}
	}
	t.amountOut = amountIn
}

func (t *innerTrade) Path() []types.Asset {
	return t.path
}

func (t *innerTrade) Pools() []types.Pool {
	return t.pools
}

func (t *innerTrade) AmountOut() *big.Int {
	return t.amountOut
}
