package consumer

type Scaler struct {
}

func NewScaler(brokerURL string) *Scaler {
	return &Scaler{}
}

func (s Scaler) ReplicaSize() int32 {
	return 2
}
