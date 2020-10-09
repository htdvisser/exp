package router

import (
	"htdvisser.dev/exp/rjs/types"
)

const TokenURI = "/api/v1/router/token"

type TokenRequestItem struct {
	Router   types.ID6 `json:"router"`
	FlavorID string    `json:"flavorid"`
}

type TokenRequest struct {
	OwnerID types.ID6 `json:"ownerid"`
	TokenRequestItem
}

type TokenBulkRequest struct {
	OwnerID types.ID6          `json:"ownerid"`
	Routers []TokenRequestItem `json:"routers"`
}

type TokenResponseItem struct {
	Router types.ID6 `json:"router"`
	Error  string    `json:"error,omitempty"`
	Token  string    `json:"token"`
}

type TokenResponse []*TokenResponseItem
