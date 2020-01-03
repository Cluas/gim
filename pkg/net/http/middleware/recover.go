package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"github.com/Cluas/gim/pkg/log"
	"github.com/Cluas/gim/pkg/mail"
)

//Recover is middleware to recover panic error
func Recover(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var sct *statusCodeTracker
		defer func() {
			// 错误 恢复
			if err := recover(); err != nil {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = fmt.Fprint(w, fmt.Sprintf(`{"code": 344, "alert": "抱歉，出错了，我们会尽快解决", "data": null, "msg": "%s"}`, err))
				//标准日志输出
				logging(sct, r)
				// 错误日志输出
				requestID := r.Context().Value(RequestIDContextKey{})
				errStack := fmt.Sprintf("request_id: %s\n\nerr_msg: %s\n\ntrace_log: %s", requestID, err, string(debug.Stack()))
				log.Bg().Error(errStack)
				// 邮件发送
				mailMsg := &mail.Message{
					Subject:   "多人点餐500错误自动通知",
					Content:   bytes.NewBufferString("<pre>" + errStack + "</pre>"),
					To:        []string{"yuyf@luedongtech.com", "huwl@luedongtech.com"},
					Extension: nil,
				}
				_ = mail.AsyncSend(mailMsg, func(err error) {
					if err != nil {
						log.Bg().Error("邮件发送错误", zap.Error(err))
					}
				})

			}
		}()
		sct = &statusCodeTracker{ResponseWriter: w}
		fn(sct.wrappedResponseWriter(), r, p)

	}
}
