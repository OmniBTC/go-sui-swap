package cetus

import (
	"context"
	"encoding/json"

	client "github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/move_types"
	suitypes "github.com/coming-chat/go-sui/v2/types"
	"github.com/omnibtc/go-sui-swap/types"
	"github.com/omnibtc/go-sui-swap/util"
)

type CetusPoolConfig struct {
	CreatePoolEventPackage string
	PoolCacheExpireSec     uint64
	Quoter                 types.Quoter
}

type createEventPoolProvider struct {
	c         *client.Client
	config    *CetusPoolConfig
	poolCache *types.PoolCache
}

type poolCreateEvent struct {
	CoinTypeA   string `json:"coin_type_a"`
	CoinTypeB   string `json:"coin_type_b"`
	PoolId      string `json:"pool_id"`
	TickSpacing int    `json:"tick_spacing"`
}

func NewCetusPoolProvider(c *client.Client, config *CetusPoolConfig) types.Provider {
	return &createEventPoolProvider{
		c:         c,
		config:    config,
		poolCache: types.NewPoolCache(config.PoolCacheExpireSec),
	}
}

func (p *createEventPoolProvider) Pools(ctx context.Context) ([]types.Pool, error) {
	pools := p.poolCache.GetPoolCache()
	if pools != nil {
		return pools, nil
	}
	pools, err := p.fetchPoolByEvents(ctx)
	if err != nil {
		p.poolCache.SetPoolsCache(pools)
	}
	return pools, err
}

func (p *createEventPoolProvider) fetchPoolByEvents(ctx context.Context) ([]types.Pool, error) {
	pageSize := uint(2000)
	hasMore := true
	events := make([]suitypes.SuiEvent, 0)
	var cursor *suitypes.EventId
	moveEventType := p.config.CreatePoolEventPackage + "::factory::CreatePoolEvent"
	for hasMore {
		data, err := p.c.QueryEvents(ctx, suitypes.EventFilter{
			MoveEventType: &moveEventType,
		}, cursor, &pageSize, false)
		if err != nil {
			return nil, err
		}
		cursor = data.NextCursor
		hasMore = data.HasNextPage
		events = append(events, data.Data...)
	}

	poolDetails := make([]types.Pool, 0)
	for _, event := range events {
		var pEvent poolCreateEvent
		data, err := json.Marshal(event.ParsedJson)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &pEvent)
		if err != nil {
			return nil, err
		}
		poolDetails = append(poolDetails, &cetusPool{
			poolAddress:  util.ShortCoinTypeWithPrefix(pEvent.PoolId),
			coinAAddress: util.ShortCoinTypeWithPrefix(pEvent.CoinTypeA),
			coinBAddress: util.ShortCoinTypeWithPrefix(pEvent.CoinTypeB),
			config:       p.config,
		})
	}

	return p.filterPausePool(ctx, poolDetails)
}

func (p *createEventPoolProvider) filterPausePool(ctx context.Context, pools []types.Pool) ([]types.Pool, error) {
	resPoolDetails := make([]types.Pool, 0)
	for i := 0; i < len(pools); i += 50 {
		r := i + 50
		if r > len(pools) {
			r = len(pools)
		}
		ps := pools[i:r]
		objectIds := []move_types.AccountAddress{}
		for _, p := range ps {
			objId, err := move_types.NewAccountAddressHex(p.Address())
			if err != nil {
				return nil, err
			}
			objectIds = append(objectIds, *objId)
		}

		objectInfos, err := p.c.MultiGetObjects(context.Background(), objectIds, &suitypes.SuiObjectDataOptions{
			ShowType:    true,
			ShowContent: true,
			ShowOwner:   true,
			ShowDisplay: true,
		})
		if err != nil {
			return nil, err
		}

		for _, poolObject := range objectInfos {
			structTag, err := util.ParseMoveStructTag(*poolObject.Data.Type)
			if err != nil {
				continue
			}
			if len(structTag.TypeParams) != 2 {
				continue
			}
			if nil == poolObject.Data ||
				nil == poolObject.Data.Content ||
				nil == poolObject.Data.Content.Data.MoveObject ||
				nil == poolObject.Data.Content.Data.MoveObject.Fields {
				continue
			}
			fieldMap, ok := poolObject.Data.Content.Data.MoveObject.Fields.(map[string]interface{})
			if !ok {
				continue
			}
			isPause, ok := fieldMap["is_pause"].(bool)
			if !ok {
				continue
			}
			if isPause {
				continue
			}

			resPoolDetails = append(resPoolDetails, &cetusPool{
				poolAddress: poolObject.Data.ObjectId.ShortString(),
				// Type:         *poolObject.Data.Type,
				coinAAddress: util.GetTypeTagFullName(structTag.TypeParams[0]),
				coinBAddress: util.GetTypeTagFullName(structTag.TypeParams[1]),
				config:       p.config,
			})
		}
	}

	return resPoolDetails, nil
}
