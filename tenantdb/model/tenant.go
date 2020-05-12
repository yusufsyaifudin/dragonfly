package model

import "time"

type Tenant struct {
	ID           string
	Name         string
	ConnectionId string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Tenants []*Tenant
