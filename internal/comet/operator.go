package comet

import "context"

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
	Connect(context.Context, *ConnectArg) (string, error)
	Disconnect(context.Context, *DisconnectArg) error
}

// DefaultOperator is default operator
type DefaultOperator struct{}

// Connect is func to Connect
func (operator *DefaultOperator) Connect(ctx context.Context, c *ConnectArg) (uid string, err error) {
	uid, err = connect(ctx, c)
	return
}

// Disconnect is func to Disconnect
func (operator *DefaultOperator) Disconnect(ctx context.Context, d *DisconnectArg) (err error) {
	if err = disconnect(ctx, d); err != nil {
		return
	}
	return
}
