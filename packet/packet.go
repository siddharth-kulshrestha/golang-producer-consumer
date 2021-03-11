package packet

type Packet struct {
	Data interface{}
	Err  error
}

func NewData(data interface{}, err error) Packet {
	return Packet{
		Data: data,
		Err:  err,
	}
}
