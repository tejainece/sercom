package comm

import (
	"github.com/tarm/serial"
	//"log"
	"time"
)

type (
	//Port is a serial port communicator
	Port struct {
		//config is the serial port configuration
		config serial.Config

		//port is the serial port
		port *serial.Port

		//Incomming buffer
		inBuf []byte

		//Outgoing buffer
		outBuf []byte

		//record determines if data must be recorded
		record bool
	}
)

//MakePort creates the port from given configuration
func MakePort(aConfig serial.Config) *Port {
	return &Port{
		config: aConfig,
	}
}

//IsRecording returns if the data is being recorded
func (meCom Port) IsRecording() bool {
	return meCom.record
}

//StartRecording starts recording. It clears the record buffer.
func (meCom *Port) StartRecording() {
	meCom.ClearRecordingBuf()
	meCom.record = true
}

//StopRecording stops recording. It doesn't clear the record buffer. Data
//in the buffer can be obtained later using GetDebugInBuf.
func (meCom *Port) StopRecording() {
	meCom.record = false
}

//ClearRecordingBuf clears recorded data
func (meCom *Port) ClearRecordingBuf() {
	meCom.inBuf = []byte{}
}

//GetRecordInBuf returns recorded in buffer
func (meCom *Port) GetRecordInBuf() []byte {
	return meCom.inBuf
}

//Open opens the serial port
func (meCom *Port) Open() error {
	if meCom.port != nil {
		return ErrAttached
	}

	var lErr error
	meCom.port, lErr = serial.OpenPort(&meCom.config)

	return lErr
}

//Close closes the serial port
func (meCom *Port) Close() error {
	if meCom.port == nil {
		return ErrNotAttached
	}

	lRet := meCom.port.Close()
	meCom.port = nil

	return lRet
}

//SendStr transmits the provided string over the serial port
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

//Send transmits the provided string over the serial port
func (meCom *Port) Send(aData []byte) error {
	if meCom.port == nil {
		return ErrNotAttached
	}

	lWritten, lErr := meCom.port.Write(aData)

	if lErr != nil {
		return lErr
	}

	if lWritten != len(aData) {
		return ErrWriteFail
	}

	return nil
}

//ReadMatchStr checks if the provided string is in the receive buffer
func (meCom *Port) ReadMatchStr(aStr string, aTmoMs int) error {
	if meCom.port == nil {
		return ErrNotAttached
	}

	lTmo := aTmoMs

	for _, cChar := range aStr {
		var bErr error
		var bByte byte

		bByte, lTmo, bErr = meCom.readByte(lTmo)

		if bErr != nil {
			return bErr
		}

		//log.Printf("%c %c", cChar, lDummy[0])

		if cChar != rune(bByte) {
			return ErrNoMatch
		}
	}

	return nil
}

//ReadBytes reads bytes from serial port receive buffer
func (meCom *Port) ReadBytes(aTmoMs int) ([]byte, error) {
	var lRet []byte
	var lRetErr error

	if meCom.port == nil {
		return lRet, ErrNotAttached
	}

	cIdx := 0

mainLoop:
	for {
		lBytes := make([]byte, 100, 100)
		bRecNum, bErr := meCom.port.Read(lBytes[:])

		if bErr != nil {
			if bErr != ErrRxTimeout {
				lRetErr = bErr
				break mainLoop
			}
		} else if bRecNum != 0 {
			lRet = append(lRet, lBytes[:bRecNum]...)
		}

		cIdx++

		if cIdx > aTmoMs {
			break
		}

		time.Sleep(time.Microsecond)
	}

	if meCom.record {
		meCom.inBuf = append(meCom.inBuf, lRet...)
	}

	return lRet, lRetErr
}

//ReadStr reads string from serial port receive buffer
func (meCom *Port) ReadStr(aTmoMs int) (string, error) {
	lString, lRetErr := meCom.ReadBytes(aTmoMs)

	return string(lString), lRetErr
}

//ReadLineStr reads a line from the serial port receive buffer
func (meCom *Port) ReadLineStr(aTmoMs int) (string, error) {
	lBytes, lErr := meCom.ReadLine(aTmoMs)

	return string(lBytes), lErr
}

//ReadLine reads a line from the serial port receive buffer
func (meCom *Port) ReadLine(aTmoMs int) ([]byte, error) {
	if meCom.port == nil {
		return nil, ErrNotAttached
	}

	var bRet []byte

	lTmo := aTmoMs

	for {
		var bErr error
		var bByte byte
		bByte, lTmo, bErr = meCom.readByte(lTmo)

		if bErr != nil {
			return bRet, bErr
		}

		if bByte == 10 {
			if len(bRet) == 0 {
				return bRet, ErrRxNoBody
			}

			bPrev := bRet[len(bRet)-1]

			if bPrev == 13 {
				bRet = bRet[:len(bRet)-1]
				break
			} else {
				return bRet, ErrRxNoCR
			}
		}

		//log.Printf("%c", rune(bByte))

		bRet = append(bRet, bByte)
	}

	return bRet, nil
}

//ReadLen reads string of given length from the serial port receive buffer
func (meCom *Port) ReadLen(aLen, aTmoMs int) ([]byte, error) {
	lRet := make([]byte, 0, aLen)
	var lRetErr error

	if meCom.port == nil {
		return lRet, ErrNotAttached
	}

	cIdx := 0
	bDummy := [1]byte{}

mainLoop:
	for cPos := 0; cPos < aLen; cPos++ {
		for {
			bRecNum, bErr := meCom.port.Read(bDummy[:])

			if bErr != nil {
				if bErr != ErrRxTimeout {
					lRetErr = bErr
					break mainLoop
				}
			} else if bRecNum == 1 {
				lRet = append(lRet, bDummy[0])
				break
			}

			cIdx++

			if cIdx > aTmoMs {
				lRetErr = ErrRxTimeout
				break mainLoop
			}

			time.Sleep(time.Microsecond)
		}
	}

	if meCom.record {
		meCom.inBuf = append(meCom.inBuf, lRet...)
	}

	return lRet, lRetErr
}

func (meCom *Port) ReadLenStr(aLen, aTmoMs int) (string, error) {
	lBytes, lErr := meCom.ReadLen(aLen, aTmoMs)

	return string(lBytes), lErr
}

//ReadSep reads string from the serial port receive buffer until it hits a seperator
func (meCom *Port) ReadSep(aSep rune, aTmoMs int) ([]byte, error) {
	if meCom.port == nil {
		return nil, ErrNotAttached
	}

	var bRet []byte

	lTmo := aTmoMs

	for {
		var bErr error
		var bByte byte
		bByte, lTmo, bErr = meCom.readByte(lTmo)

		if bErr != nil {
			return bRet, bErr
		}

		if bByte == byte(aSep) {
			return bRet, nil
		}

		//log.Printf("%c", rune(bDummy[0]))

		bRet = append(bRet, bByte)
	}
}

//ReadSep reads string from the serial port receive buffer until it hits a seperator
func (meCom *Port) ReadSepStr(aSep rune, aTmoMs int) (string, error) {
	lBytes, lErr := meCom.ReadSep(aSep, aTmoMs)

	return string(lBytes), lErr
}

func (meCom *Port) readByte(aTmoIn int) (byte, int, error) {
	if meCom.port == nil {
		return 0, aTmoIn, ErrNotAttached
	}

	bDummy := [1]byte{}

	for {
		bNum, bErr := meCom.port.Read(bDummy[:])

		if bErr != nil {
			return 0, aTmoIn, bErr
		}

		if bNum != 0 {
			break
		}

		if aTmoIn == 0 {
			return 0, aTmoIn, ErrRxTimeout
		}

		aTmoIn--

		time.Sleep(time.Microsecond)
	}

	if meCom.record {
		meCom.inBuf = append(meCom.inBuf, bDummy[0])
	}

	return bDummy[0], aTmoIn, nil
}
