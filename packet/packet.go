package packet

type Packet struct {
	Data interface{}
	Err  error
}

func New(data interface{}, err error) Packet {
	return Packet{
		Data: data,
		Err:  err,
	}
}
