package comm

import (
	"github.com/tarm/serial"
	"time"
)

type (
	//Port is a serial port communicator
	Port struct {
		config serial.Config
		port   *serial.Port

		//Incomming buffer
		inBuf string

		//Outgoing buffer
		outBuf string

		record bool
	}
)

func MakePort(aConfig serial.Config) *Port {
	return &Port{
		config: aConfig,
	}
}

func (meCom Port) IsRecording() bool {
	return meCom.record
}

func (meCom *Port) StartRecording() {
	meCom.ClearDebug()
	meCom.record = true
}

func (meCom *Port) StopRecording() {
	meCom.record = false
}

func (meCom *Port) ClearDebug() {
	meCom.inBuf = ""
}

func (meCom *Port) GetDebugInBuf() string {
	return meCom.inBuf
}

func (meCom *Port) Open() error {
	if meCom.port != nil {
		return ErrAttached
	}

	var lErr error
	meCom.port, lErr = serial.OpenPort(&meCom.config)

	return lErr
}

func (meCom *Port) Close() error {
	if meCom.port == nil {
		return ErrNotAttached
	}

	lRet := meCom.port.Close()
	meCom.port = nil

	return lRet
}

func (meCom *Port) SendStr(aStr string) error {
	if meCom.port == nil {
		return ErrNotAttached
	}

	lWritten, lErr := meCom.port.Write([]byte(aStr))

	if lErr != nil {
		return lErr
	}

	if lWritten != len(aStr) {
		return ErrWriteFail
	}

	return nil
}

func (meCom *Port) ReadMatchStr(aStr string, aTmoMs int) error {
	if meCom.port == nil {
		return ErrNotAttached
	}

	lDummy := [1]byte{}
	cIdx := 0
	for _, cChar := range aStr {
		for {
			bNum, bErr := meCom.port.Read(lDummy[:])

			if bErr != nil {
				return bErr
			}

			if bNum != 0 {
				break
			}

			cIdx++

			if cIdx > aTmoMs {
				return ErrRxTimeout
			}

			time.Sleep(time.Millisecond)
		}

		if meCom.record {
			meCom.inBuf += string(lDummy[0])
		}
		//log.Printf("%c %c", cChar, lDummy[0])

		if cChar != rune(lDummy[0]) {
			return ErrNoMatch
		}
	}

	return nil
}

func (meCom *Port) ReadStr(aTmoMs int) (string, error) {
	if meCom.port == nil {
		return "", ErrNotAttached
	}

	lRet := make([]byte, 0)

	lLen := 0

	for cIdx := 0; ; {
		lBytes := make([]byte, 100, 100)
		bRecNum, bErr := meCom.port.Read(lBytes[:])

		if bErr != nil {
			if bErr != ErrRxTimeout {
				return "", bErr
			} else {
				//TODO
			}
		} else {
			if bRecNum != 0 {
				lLen += bRecNum
				lRet = append(lRet, lBytes[:bRecNum]...)
			}
		}

		cIdx++

		if cIdx > aTmoMs {
			break
		}

		time.Sleep(time.Microsecond)
	}

	if meCom.record {
		meCom.inBuf += string(lRet[:lLen])
	}

	return string(lRet[:lLen]), nil
}

func (meCom *Port) ReadLineStr(aTmoMs int) (string, error) {
	if meCom.port == nil {
		return "", ErrNotAttached
	}

	lBytes, lErr := meCom.ReadLine(aTmoMs)

	if lErr != nil {
		return "", lErr
	}

	return string(lBytes), nil
}

// ReadLine reads a line from the given serial port
func (meCom *Port) ReadLine(aTmoMs int) ([]byte, error) {
	if meCom.port == nil {
		return nil, ErrNotAttached
	}

	bRet := make([]byte, 0)

	cIdx := 0
	for {
		bDummy := [1]byte{}
		for {
			bNum, bErr := meCom.port.Read(bDummy[:])

			if bErr != nil {
				return bRet, bErr
			}

			if bNum != 0 {
				break
			}

			cIdx++

			if cIdx > aTmoMs {
				return bRet, ErrRxTimeout
			}

			time.Sleep(time.Microsecond)
		}

		if meCom.record {
			meCom.inBuf += string(bDummy[0])
		}

		if bDummy[0] == 10 {
			if len(bDummy) == 0 {
				bRet = append(bRet, bDummy[0])
				return bRet, ErrRxNoBody
			}

			bPrev := bRet[len(bRet)-1]
			bRet = append(bRet, bDummy[0])

			if bPrev == 13 {
				bRet = bRet[:len(bRet)-2]
				break
			} else {
				return bRet, ErrRxNoCR
			}
		}

		bRet = append(bRet, bDummy[0])

	}

	return bRet, nil
}
