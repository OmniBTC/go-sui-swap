package suiswap

import (
	"math/big"
	"sync"

	"github.com/omnibtc/go-sui-swap/types"
)

type innerTrade struct {
	pools      []types.Pool
	path       []types.Asset
	isA2B      []bool
	amountOuts []*big.Int

	wg *sync.WaitGroup
}

func newTrade(pools []types.Pool, path []string, isA2B []bool, amountIn *big.Int, wg *sync.WaitGroup) types.Trade {
	assets := make([]types.Asset, 0, len(path))
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
	amountOuts := make([]*big.Int, 0)
	for i, p := range t.pools {
		amountIn, err = p.GetAmountOut(p, t.path[i], t.path[i+1], amountIn)
		if err != nil {
			return
		}
		amountOuts = append(amountOuts, amountIn)
	}
	t.amountOuts = amountOuts
}

func (t *innerTrade) Path() []types.Asset {
	return t.path
}

func (t *innerTrade) Pools() []types.Pool {
	return t.pools
}

func (t *innerTrade) AmountOuts() []*big.Int {
	return t.amountOuts
}

func (t *innerTrade) A2B() []bool {
	return t.isA2B
}
