package types

import "time"

type PoolCache struct {
	pools         []Pool
	lastCacheTime uint64
	expire        uint64
}

func NewPoolCache(expireSec uint64) *PoolCache {
	if expireSec == 0 {
		expireSec = 3600
	}
	return &PoolCache{expire: expireSec}
}

func (c *PoolCache) GetPoolCache() []Pool {
	if c.lastCacheTime+uint64(c.expire) < uint64(time.Now().Unix()) {
		return nil
	}
	return c.pools
}

func (c *PoolCache) SetPoolsCache(pools []Pool) {
	c.pools = pools
	c.lastCacheTime = uint64(time.Now().Unix())
}
