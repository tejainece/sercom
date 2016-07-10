package comm

//Err contains information about the failure
type Err string

//Error returns printable string of the failure
func (meErr Err) Error() string {
	return string(meErr)
}

const (
	ErrAttached    = Err("Comm has another serial port attached!")
	ErrNotAttached = Err("Comm has no serial port attached!")
	ErrError       = Err("Unknown error. Please file a report!")
	ErrRxTimeout   = Err("Reception timed out!")
	ErrRxNoBody    = Err("Empty body!")
	ErrRxNoCR      = Err("No carriage return in trailer!")
	ErrNoMatch     = Err("Match not found!")
	ErrWriteFail   = Err("Command write failed!")
	ErrFormat      = Err("Format error!")
)
