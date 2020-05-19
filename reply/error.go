package reply

import (
	"encoding/json"
	"net/http"
	"ysf/dragonfly/server"
)

type replyError struct {
	data interface{}
}

func (s replyError) StatusCode() int {
	return http.StatusUnprocessableEntity
}

func (s replyError) Body() (data []byte, err error) {
	data, err = json.Marshal(s.data)
	return
}

func (s replyError) Header() http.Header {
	return http.Header{}
}

func (s replyError) ContentType() string {
	return server.ContentTypeJSON
}

func Error(data interface{}) server.Response {
	return &replyError{
		data: data,
	}
}
