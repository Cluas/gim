package mail

// Send is func to send mail
func Send(msg *Message) (err error) {
	err = sender.Send(msg)
	return
}

// AsyncSend is func to async send mail
func AsyncSend(msg *Message, handle func(err error)) (err error) {
	err = sender.AsyncSend(msg, handle)
	return
}
