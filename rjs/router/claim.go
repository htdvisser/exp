package router

import "htdvisser.dev/exp/rjs/types"

const ClaimURI = "/api/v1/router/claim"

type ClaimRequestItem struct {
	Router types.ID6 `json:"router"`
	Claim  string    `json:"claim"`
}

type ClaimRequest struct {
	OwnerID types.ID6 `json:"ownerid"`
	ClaimRequestItem
}

type ClaimBulkRequest struct {
	OwnerID types.ID6          `json:"ownerid"`
	Routers []ClaimRequestItem `json:"routers"`
}

type ClaimResponseItem struct {
	Router types.ID6 `json:"router"`
	Error  string    `json:"error,omitempty"`
}

type ClaimResponse []*ClaimResponseItem
