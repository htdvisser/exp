package router

import "htdvisser.dev/exp/rjs/types"

const DeleteURI = "/api/v1/router/delete"

type DeleteRequest struct {
	OwnerID types.ID6 `json:"ownerid"`
	Router  types.ID6 `json:"router"`
}

type DeleteBulkRequest struct {
	OwnerID types.ID6   `json:"ownerid"`
	Routers []types.ID6 `json:"routers"`
}

type DeleteResponseItem struct {
	Router types.ID6 `json:"router"`
	Error  string    `json:"error,omitempty"`
}

type DeleteResponse []*DeleteResponseItem
