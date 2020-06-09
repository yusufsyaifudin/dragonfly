package processor

import (
	"context"
	"encoding/json"
	"strings"
	"ysf/dragonfly/reply"
	"ysf/dragonfly/server"

	"github.com/opentracing/opentracing-go"
)

type serverReqToInput struct {
	req server.Request
}

func (s serverReqToInput) GetData() ([]byte, error) {
	var data interface{}
	if err := s.req.Bind(&data); err != nil {
		return nil, err
	}

	return json.Marshal(data)
}

func (s serverReqToInput) GetHeader(key string) string {
	return s.req.RawRequest().Header.Get(key)
}

func (s serverReqToInput) GetQueryParam(key string) string {
	tag := strings.TrimSpace(s.req.GetQueryParam(key))
	return tag
}

func (s serverReqToInput) GetParam(key string) string {
	tag := strings.TrimSpace(s.req.GetParam(key))
	return tag
}

func ToServerHandler(service Service) server.Handler {
	return func(ctx context.Context, request server.Request) server.Response {
		span, ctx := opentracing.StartSpanFromContext(ctx, "ToServerHandler")
		defer func() {
			span.Finish()
			ctx.Done()
		}()

		out := service(ctx, &serverReqToInput{req: request})
		if out.Error != nil {
			return reply.Error(map[string]interface{}{
				"code": out.StatusCode,
				"type": out.Type,
				"error": map[string]interface{}{
					"message": out.Error.Error(),
				},
			})
		}

		return reply.Success(map[string]interface{}{
			"code": out.StatusCode,
			"type": out.Type,
			"data": out.Data,
		})
	}
}
