package gin_response_handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

type Middleware struct {
	TelegramBotId  string
	TelegramChatId string
}

func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func NewMiddleware(args Middleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		args.ginResponseMiddleware(c)
	}
}

func (args *Middleware) ginResponseMiddleware(c *gin.Context) {
	bw := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = bw
	c.Next()

	statusCode := c.Writer.Status()
	path := c.FullPath()
	method := c.Request.Method
	proto := c.Request.Proto

	reqBody, _ := ioutil.ReadAll(c.Request.Body)

	if statusCode >= http.StatusBadRequest {
		msg := fmt.Sprintf("%s%%20%s%%20%s%%0AResponse:%%20%s",
			method, path, proto, bw.body.String())

		if len(reqBody) > 0 {
			rmsg := fmt.Sprintf("%s%%0ARequest:%%20%s", msg, string(reqBody))
			msg = rmsg
		}

		args.sendTelegram(msg)
	}
}

func (args *Middleware) sendTelegram(msg string) {
	url := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage?chat_id=%s&parse_mode=Markdown&text=%s",
		args.TelegramBotId, args.TelegramChatId, msg)
	_, _ = http.Get(url)
}
