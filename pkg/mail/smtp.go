package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
)

var sender *SMTPSender

// Config is struct of sender config
type Config struct {
	Addr     string
	From     string
	Password string
	StartSSL bool
}

func init() {
	sender = &SMTPSender{
		addr: "smtp.mxhichina.com:465",
		from: mail.Address{Name: "", Address: "bug@luedongtech.com"},
		cfg: &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         "smtp.mxhichina.com",
		},
	}
	sender.auth = smtp.PlainAuth("", sender.from.Address, "Ldtech123#", "smtp.mxhichina.com")

}

// Init provide a hook to use own config
func Init(c *Config) (err error) {
	sender.addr = c.Addr
	sender.from = mail.Address{Name: "", Address: c.From}
	host, _, err := net.SplitHostPort(c.Addr)
	if err != nil {
		return
	}
	if c.StartSSL {
		sender.cfg = &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}
	}
	sender.auth = smtp.PlainAuth("", sender.from.Address, c.Password, host)
	return
}

// SMTPSender used to send mail with smtp protocol
type SMTPSender struct {
	addr string
	from mail.Address
	auth smtp.Auth
	cfg  *tls.Config
}

// Send send mail
func (s *SMTPSender) Send(msg *Message) (err error) {

	err = s.send(msg)

	return
}

// AsyncSend async send mail
func (s *SMTPSender) AsyncSend(msg *Message, handle func(err error)) error {
	go func() {
		err := s.send(msg)
		handle(err)
	}()
	return nil
}

func (s *SMTPSender) send(msg *Message) error {
	for _, addr := range msg.To {
		header := s.getHeader(msg.Subject)
		header["To"] = addr
		if msg.Extension != nil {
			for k, v := range msg.Extension {
				header[k] = v
			}
		}
		data := s.getData(header, msg.Content)
		conn, err := tls.Dial("tcp", s.addr, s.cfg)
		if err != nil {
			log.Panic(err)
		}
		host, _, err := net.SplitHostPort(s.addr)
		client, err := smtp.NewClient(conn, host)

		if err != nil {
			log.Panic(err)
		}
		if err = client.Auth(s.auth); err != nil {
			log.Panic(err)
		}

		if err = client.Mail(s.from.Address); err != nil {
			log.Panic(err)
		}

		if err = client.Rcpt(addr); err != nil {
			log.Panic(err)
		}

		w, err := client.Data()
		if err != nil {
			log.Panic(err)
		}
		_, err = w.Write(data)
		err = w.Close()
		if err != nil {
			log.Panic(err)
		}
		err = client.Quit()
		if err != nil {
			log.Panic(err)
		}
	}
	return nil
}

func (s *SMTPSender) getHeader(subject string) map[string]string {
	header := make(map[string]string)
	header["From"] = s.from.String()
	header["Subject"] = mime.QEncoding.Encode("utf-8", subject)
	header["Mime-Version"] = "1.0"
	header["Content-Type"] = "text/html;charset=utf-8"
	header["Content-Transfer-Encoding"] = "Quoted-Printable"
	return header
}

func (s *SMTPSender) getData(header map[string]string, body io.Reader) []byte {
	buf := new(bytes.Buffer)
	for k, v := range header {
		_, _ = fmt.Fprintf(buf, "%s: %s\r\n", k, v)
	}
	_, _ = fmt.Fprintf(buf, "\r\n")
	_, _ = io.Copy(buf, body)
	return buf.Bytes()
}
