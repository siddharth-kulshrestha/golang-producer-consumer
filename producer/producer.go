package producer

type Data struct {
	Data interface{}
	Err  error
}

func NewData(data interface{}, err error) Data {
	return Data{
		Data: data,
		Err:  err,
	}
}

type ProducerFunc func(args ...interface{}) (Data, bool, []interface{})

type Producer interface {
	Produce(chan Data, ProducerFunc)
	CallProducerFunc()
}
