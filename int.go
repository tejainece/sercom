package comm

import (
	//"log"
	"strconv"
)

type Radix uint8

const (
	RadixAuto = Radix(iota)
	RadixBinary
	RadixDecimal
	RadixHex
)

type BitSize uint8

const (
	BitSizeFull = BitSize(iota)
	BitSize8
	BitSize16
	BitSize32
	BitSize64
)

func RadixToInt(aRadix Radix) int {
	switch aRadix {
	case RadixAuto:
		return 0
	case RadixBinary:
		return 2
	case RadixDecimal:
		return 10
	case RadixHex:
		return 16
	default:
		return 0
	}
}

func BitSizeToInt(aBitSize BitSize) int {
	switch aBitSize {
	case BitSizeFull:
		return 0
	case BitSize8:
		return 8
	case BitSize16:
		return 16
	case BitSize32:
		return 32
	case BitSize64:
		return 64
	default:
		return 0
	}
}

//ReadInt reads int32 from decimal string
func (meCom *Port) ReadInt(aRadix Radix, aSize BitSize, aTmo int) (int64, error) {
	var lStr string
	var lErr error

	if lStr, lErr = meCom.ReadStr(aTmo); lErr != nil {
		return 0, lErr
	}

	lInt, lErr := strconv.ParseInt(lStr, RadixToInt(aRadix), BitSizeToInt(aSize))

	if lErr != nil {
		return 0, ErrFormat
	}

	return lInt, lErr
}

//ReadUint reads int32 from decimal string
func (meCom *Port) ReadUint(aRadix Radix, aSize BitSize, aTmo int) (uint64, error) {
	var lStr string
	var lErr error

	if lStr, lErr = meCom.ReadStr(aTmo); lErr != nil {
		return 0, lErr
	}

	lInt, lErr := strconv.ParseUint(lStr, RadixToInt(aRadix), BitSizeToInt(aSize))

	if lErr != nil {
		return 0, ErrFormat
	}

	return lInt, lErr
}
