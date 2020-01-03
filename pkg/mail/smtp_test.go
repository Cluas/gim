package mail_test

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/Cluas/gim/pkg/mail"
)

func ExampleSend() {
	err := mail.Send(&mail.Message{
		Subject: "测试",
		Content: bytes.NewBufferString("<h1>你好，同步测试邮件内容</h1>"),
		To:      []string{"huwl@luedongtech.com"},
		//Extension: nil,
	})
	if err != nil {
		fmt.Printf("发送错误, %v", err)
	}
	// Output:
}

func ExampleAsyncSend() {
	var wg sync.WaitGroup
	wg.Add(1)
	err := mail.AsyncSend(&mail.Message{
		Subject: "测试",
		Content: bytes.NewBufferString("<h1>你好，异步发送测试邮件内容</h1>"),
		To:      []string{"huwl@luedongtech.com"},
		//Extension: nil,
	},
		func(err error) {
			defer wg.Done()
			if err != nil {
				fmt.Println("发送邮件出现错误：", err)
			}
		})
	if err != nil {
		fmt.Printf("发送错误, %v", err)
	}
	wg.Wait()
	// Output:
}
