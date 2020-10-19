package router

import "htdvisser.dev/exp/rjs/types"

const AddURI = "/api/v1/router/add"

type AddRequestItem struct {
	Router   types.ID6 `json:"router"`
	FlavorID string    `json:"flavorid"`
	Token    []byte    `json:"token"`
}

type AddRequest struct {
	OwnerID types.ID6 `json:"ownerid"`
	AddRequestItem
}

type AddBulkRequest struct {
	OwnerID types.ID6        `json:"ownerid"`
	Routers []AddRequestItem `json:"routers"`
}

type AddResponseItem struct {
	Router types.ID6 `json:"router"`
	Error  string    `json:"error,omitempty"`
}

type AddResponse []*AddResponseItem
