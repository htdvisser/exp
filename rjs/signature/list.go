package signature

import "htdvisser.dev/exp/rjs/types"

const ListURI = "/api/v1/signature/list"

type ListRequest struct {
	OwnerID types.ID6 `json:"ownerid"`
}

type ListResponseItem struct {
	SigCRC int32 `json:"sigcrc"`
	FWCRC  int32 `json:"fwcrc"`
	KeyCRC int32 `json:"keycrc"`
}

type ListResponse []*ListResponseItem
