package suiswap

import (
	"context"
	"math/big"
	"sort"
	"sync"

	"github.com/omnibtc/go-sui-swap/types"
	"github.com/omnibtc/go-sui-swap/util"
)

type BestTradeOptions struct {
	LimitTradeCount int
}

type App struct {
	providers []types.Provider
}

func NewApp() *App {
	return &App{}
}

func (a *App) RegisterProvider(provider types.Provider) {
	a.providers = append(a.providers, provider)
}

func (a *App) Pools(ctx context.Context) ([]types.Pool, error) {
	ps := make([]types.Pool, 0, 8)
	var lastErr error
	for _, p := range a.providers {
		tps, err := p.Pools(ctx)
		if err != nil {
			lastErr = err
			continue
		}
		ps = append(ps, tps...)
	}
	if len(ps) == 0 {
		return ps, lastErr
	}
	return ps, nil
}

func (a *App) BestTradeExactIn(ctx context.Context, pools []types.Pool, coinIn, coinOut types.Asset, amountIn *big.Int, options BestTradeOptions) ([]types.Trade, error) {
	if options.LimitTradeCount == 0 {
		options.LimitTradeCount = 3
	}

	allPools, err := a.Pools(ctx)
	if err != nil {
		return nil, err
	}
	if len(allPools) == 0 {
		return []types.Trade{}, nil
	}

	trades := a.TokenRouter(allPools, coinIn.Address(), coinOut.Address(), amountIn)
	if len(trades) == 0 {
		return []types.Trade{}, nil
	}

	sort.Slice(trades, func(i, j int) bool {
		return trades[i].AmountOut().Cmp(trades[j].AmountOut()) >= 0
	})

	if len(trades) > options.LimitTradeCount {
		return trades[:options.LimitTradeCount], nil
	}

	return trades, nil
}

func (a *App) TokenRouter(pools []types.Pool, coinIn string, coinOut string, amountIn *big.Int) []types.Trade {
	trades := make([]types.Trade, 0)
	wg := &sync.WaitGroup{}
	coin2pools := make(map[string][]types.Pool)

	// one step
	for _, pool := range pools {
		coin2pools[pool.CoinA()] = append(coin2pools[pool.CoinA()], pool)
		coin2pools[pool.CoinB()] = append(coin2pools[pool.CoinB()], pool)

		isPool, isA2b := a.isPoolMatch(pool, coinIn, coinOut)
		if isPool {
			trades = append(trades, newTrade([]types.Pool{pool}, []string{coinIn, coinOut}, []bool{isA2b}, amountIn, wg))
		}
	}

	// two step
	for _, coinInPool := range coin2pools[coinIn] {
		middleCoin := coinInPool.CoinB()
		firstPoolIsA2B := true
		if util.EqualSuiCoinAddress(coinIn, middleCoin) {
			firstPoolIsA2B = false
			middleCoin = coinInPool.CoinA()
		}
		if util.EqualSuiCoinAddress(middleCoin, coinOut) {
			continue
		}

		for _, coinOutPool := range coin2pools[middleCoin] {
			if util.EqualSuiCoinAddress(coinOutPool.CoinA(), coinOut) ||
				util.EqualSuiCoinAddress(coinOutPool.CoinB(), coinOut) {
				trades = append(trades, newTrade(
					[]types.Pool{coinInPool, coinOutPool},
					[]string{coinIn, middleCoin, coinOut},
					[]bool{firstPoolIsA2B, util.EqualSuiCoinAddress(coinOutPool.CoinB(), coinOut)},
					amountIn,
					wg,
				))
			}
		}
	}
	wg.Wait()

	// filter nil outAmount trade
	validTrade := make([]types.Trade, 0, len(trades))
	for _, trade := range trades {
		if trade.AmountOut() != nil {
			validTrade = append(validTrade, trade)
		}
	}
	return validTrade
}

func (a *App) isPoolMatch(pool types.Pool, coinA string, coinB string) (isPool bool, isA2b bool) {
	if util.EqualSuiCoinAddress(coinA, pool.CoinA()) {
		if util.EqualSuiCoinAddress(coinB, pool.CoinB()) {
			return true, true
		} else {
			return false, false
		}
	} else if util.EqualSuiCoinAddress(coinA, pool.CoinB()) {
		if util.EqualSuiCoinAddress(coinB, pool.CoinA()) {
			return true, false
		} else {
			return false, false
		}
	} else {
		return false, false
	}
}
