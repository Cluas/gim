package mail

import "io"

// Message is struct used to send msg
type Message struct {
	Subject   string
	Content   io.Reader         // support html content
	To        []string          // to address string
	Extension map[string]string // message extension
}

// Sender is interface for sender
type Sender interface {
	// Send send mail
	Send(msg *Message) error
	// AsyncSend async send mail need callback func
	AsyncSend(msg *Message, handle func(err error)) error
}
