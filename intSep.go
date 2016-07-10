package comm

import (
	//"log"
	"strconv"
)

//ReadInt reads int32 from decimal string
func (meCom *Port) ReadIntSep(aRadix Radix, aSize BitSize, aSep rune, aTmo int) (int64, error) {
	var lStr string
	var lErr error

	if lStr, lErr = meCom.ReadSepStr(aSep, aTmo); lErr != nil {
		return 0, lErr
	}

	lInt, lErr := strconv.ParseInt(lStr, RadixToInt(aRadix), BitSizeToInt(aSize))

	if lErr != nil {
		return 0, ErrFormat
	}

	return lInt, lErr
}

//ReadUint reads int32 from decimal string
func (meCom *Port) ReadUintSep(aRadix Radix, aSize BitSize, aSep rune, aTmo int) (uint64, error) {
	var lStr string
	var lErr error

	if lStr, lErr = meCom.ReadSepStr(aSep, aTmo); lErr != nil {
		return 0, lErr
	}

	lInt, lErr := strconv.ParseUint(lStr, RadixToInt(aRadix), BitSizeToInt(aSize))

	if lErr != nil {
		return 0, ErrFormat
	}

	return lInt, lErr
}
