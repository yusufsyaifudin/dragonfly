package userhandler

import (
	"ysf/dragonfly/dependency"
	"ysf/dragonfly/server"
)

type handler struct {
	dep dependency.Service
}

func Routes(dep dependency.Service) []*server.Route {
	h := handler{
		dep: dep,
	}
	return []*server.Route{
		{
			Path:       "/users/:tenant_id",
			Method:     "GET",
			Handler:    h.registerUser,
			Middleware: nil,
		},
	}
}
