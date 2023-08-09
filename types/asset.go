package types

type defaultAsset struct {
	address string
}

func NewAsset(address string) Asset {
	return &defaultAsset{address: address}
}

func (a *defaultAsset) Address() string {
	return a.address
}
