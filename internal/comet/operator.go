package comet

// ConnectArg is rpc connect arg
type ConnectArg struct {
	Auth     string
	RoomID   string
	ServerID string
}

// DisconnectArg is rpc disconnect arg
type DisconnectArg struct {
	RoomID string
	UID    string
}

// Operator is interface for operation
type Operator interface {
	Connect(*ConnectArg) (string, error)
	Disconnect(*DisconnectArg) error
}

// DefaultOperator is default operator
type DefaultOperator struct{}

// Connect is func to Connect
func (operator *DefaultOperator) Connect(c *ConnectArg) (uid string, err error) {
	uid, err = connect(c)
	return
}

// Disconnect is func to Disconnect
func (operator *DefaultOperator) Disconnect(d *DisconnectArg) (err error) {
	if err = disconnect(d); err != nil {
		return
	}
	return
}
