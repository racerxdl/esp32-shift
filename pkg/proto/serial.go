package shift

type IOLevel uint8

const (
	HIGH IOLevel = 1
	LOW  IOLevel = 0
)

func (io IOLevel) String() string {
	if io == LOW {
		return "LOW"
	}

	return "HIGH"
}

type BoardStatus []IOLevel

func (s BoardStatus) String() string {
	d := "BoardStatus["

	for i, v := range s {
		d += v.String()
		if i < len(s)-1 {
			d += ", "
		}
	}

	return d + "]"
}

type ESPShift interface {
	Reset() error
	HealthCheck() error
	SetPin(pin uint8, val IOLevel) error
	SetByte(byteNum uint8, val uint8) error
	Status() (BoardStatus, error)
	Close() error
	ReadAllLines() []string
}
