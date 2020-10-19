package firmware

import "htdvisser.dev/exp/rjs/types"

const ListURI = "/api/v1/firmware/list"

type ListRequest struct {
	OwnerID types.ID6 `json:"ownerid"`
}

type ListResponseItem struct {
	Version string `json:"version"`
	FWCRC   int32  `json:"fwcrc"`
	Model   string `json:"model"`
}

type ListResponse []*ListResponseItem
