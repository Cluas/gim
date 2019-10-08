package comet

type ConnectArg struct {
	Auth     string
	RoomID   int32
	ServerID int8
}

type DisconnectArg struct {
	RoomID int32
	UID    string
}

type Operator interface {
	Connect(*ConnectArg) (string, error)
	Disconnect(*DisconnectArg) error
}

type DefaultOperator struct{}

func (operator *DefaultOperator) Connect(c *ConnectArg) (uid string, err error) {
	uid, err = connect(c)
	return
}

func (operator *DefaultOperator) Disconnect(d *DisconnectArg) (err error) {
	if err = disconnect(d); err != nil {
		return
	}
	return
}
