package userhandler

import (
	"context"
	"fmt"
	"ysf/dragonfly/processor"
)

func (h handler) registerUser(ctx context.Context, input processor.Input) processor.Output {
	tenantID := input.GetParam("tenant_id")

	tenantDB := h.dep.TenantDB()
	tenant, _, err := tenantDB.CreateTenant(ctx, tenantID, tenantID)
	if err != nil {
		return processor.Output{
			Error:      err,
			Type:       "",
			StatusCode: "",
			Data:       nil,
		}
	}

	immigration, err := tenantDB.GetTenantImmigration(ctx, tenantID)
	if err != nil {
		return processor.Output{
			Error:      fmt.Errorf("immigration error: %s", err.Error()),
			Type:       "",
			StatusCode: "",
			Data:       nil,
		}
	}

	err = immigration.Sync(ctx)
	if err != nil {
		return processor.Output{
			Error:      fmt.Errorf("sync error: %s", err.Error()),
			Type:       "",
			StatusCode: "",
			Data:       nil,
		}
	}

	return processor.Output{
		Error:      nil,
		Type:       "Tenant",
		StatusCode: "",
		Data:       tenant,
	}
}
