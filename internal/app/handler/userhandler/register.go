package userhandler

import (
	"context"
	"fmt"
	"ysf/dragonfly/reply"
	"ysf/dragonfly/server"
)

func (h *handler) registerUser(ctx context.Context, req server.Request) server.Response {
	tenantID := req.GetParam("tenant_id")

	tenantDB := h.dep.TenantDB()
	tenantConnection, err := tenantDB.CreateTenant(ctx, tenantID, tenantID)
	if err != nil {
		return reply.Success(map[string]interface{}{
			"error": fmt.Sprintf("error creating db connection: %s", err.Error()),
		})
	}

	immigration, err := tenantDB.GetTenantImmigration(ctx, tenantID)
	if err != nil {
		return reply.Success(map[string]interface{}{
			"error": fmt.Sprintf("immigration error: %s", err.Error()),
		})
	}

	err = immigration.Sync(ctx)
	if err != nil {
		return reply.Success(map[string]interface{}{
			"error": fmt.Sprintf("sync error: %s", err.Error()),
		})
	}

	return reply.Success(server.ReplyStructure{
		Type: "Tenant",
		Data: tenantConnection.TenantInfo(),
	})
}
