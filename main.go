package gin_response_handler

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

type responseBody struct {
	message string `json:"message,omitempty"`
	code    string `json:"code,omitempty"`
}

func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func GinResponseMiddleware(c *gin.Context) {
	bw := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = bw
	c.Next()

	statusCode := c.Writer.Status()
	path := c.FullPath()
	method := c.Request.Method
	proto := c.Request.Proto

	reqBody, _ := ioutil.ReadAll(c.Request.Body)

	if statusCode >= 400 {
		msg := fmt.Sprintf("%s %s %s\nResponse: %s", method, path, proto, bw.body.String())

		if len(reqBody) > 0 {
			msg = fmt.Sprintf("%s\nRequest: %s", msg, string(reqBody))
		}

		fmt.Println(msg)
	}
}
