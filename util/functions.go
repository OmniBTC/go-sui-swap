package util

import "strings"

func ShortCoinTypeWithPrefix(address string) string {
	return "0x" + strings.TrimLeft(address, "x0")
}

func FixHex(h string) string {
	h = strings.TrimPrefix(h, "0x")
	if len(h)%2 == 0 {
		return h
	}
	return "0" + h
}

func EqualSuiCoinAddress(x, y string) bool {
	var (
		ix   = 0
		iy   = 0
		c    rune
		lenx = len(x)
		leny = len(y)
	)
	for ix, c = range x {
		if c == 'x' || c == '0' {
			continue
		} else {
			break
		}
	}
	for iy, c = range y {
		if c == 'x' || c == '0' {
			continue
		} else {
			break
		}
	}
	for ix < lenx && iy < leny {
		if x[ix] != y[iy] {
			break
		}
		ix++
		iy++
	}
	return ix == lenx && iy == leny
}
