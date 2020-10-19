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

func (res ListResponse) Len() int           { return len(res) }
func (res ListResponse) Swap(i, j int)      { res[i], res[j] = res[j], res[i] }
func (res ListResponse) Less(i, j int) bool { return res[i].SigCRC < res[j].SigCRC }
