package comm

import (
	//"log"
	"strconv"
)

//ReadInt reads int32 from decimal string
func (meCom *Port) ReadIntLen(aRadix Radix, aSize BitSize, aLen, aTmo int) (int64, error) {
	var lStr string
	var lErr error

	if lStr, lErr = meCom.ReadLenStr(aLen, aTmo); lErr != nil {
		return 0, lErr
	}

	lInt, lErr := strconv.ParseInt(lStr, RadixToInt(aRadix), BitSizeToInt(aSize))

	if lErr != nil {
		return 0, ErrFormat
	}

	return lInt, lErr
}

//ReadUint reads int32 from decimal string
func (meCom *Port) ReadUintLen(aRadix Radix, aSize BitSize, aLen, aTmo int) (uint64, error) {
	var lStr string
	var lErr error

	if lStr, lErr = meCom.ReadLenStr(aLen, aTmo); lErr != nil {
		return 0, lErr
	}

	lInt, lErr := strconv.ParseUint(lStr, RadixToInt(aRadix), BitSizeToInt(aSize))

	if lErr != nil {
		return 0, ErrFormat
	}

	return lInt, lErr
}
