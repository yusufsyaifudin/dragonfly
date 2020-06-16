package reply

import (
	"encoding/json"
	"net/http"

	"github.com/yusufsyaifudin/dragonfly/server"
)

type success struct {
	data interface{}
}

func (s success) StatusCode() int {
	return http.StatusOK
}

func (s success) Body() (data []byte, err error) {
	data, err = json.Marshal(s.data)
	return
}

func (s success) Header() http.Header {
	return http.Header{}
}

func (s success) ContentType() string {
	return server.ContentTypeJSON
}

func Success(data interface{}) server.Response {
	return &success{
		data: data,
	}
}
