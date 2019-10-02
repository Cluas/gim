package comet

type ConnArg struct {
	Auth     string
	RoomId   int32
	ServerId int8
}

type DisConnArg struct {
	RoomID int32
	Uid    string
}

type Operator interface {
	Connect(*ConnArg) (string, error)
	DisConnect(*DisConnArg) error
}

type DefaultOperator struct {
}

func (operator *DefaultOperator) Connect(connArg *ConnArg) (uid string, err error) {

	return
}

func (operator *DefaultOperator) DisConnect(dArg *DisConnArg) (err error) {

	return
}
