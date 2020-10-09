package router

import (
	"htdvisser.dev/exp/rjs/types"
)

const ListURI = "/api/v1/router/list"

type ListRequest struct {
	OwnerID types.ID6 `json:"ownerid"`
}

type ListResponse []types.ID6

func (res ListResponse) Len() int           { return len(res) }
func (res ListResponse) Swap(i, j int)      { res[i], res[j] = res[j], res[i] }
func (res ListResponse) Less(i, j int) bool { return res[i].EUIString() < res[j].EUIString() }
