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

func (res ListResponse) Len() int      { return len(res) }
func (res ListResponse) Swap(i, j int) { res[i], res[j] = res[j], res[i] }
func (res ListResponse) Less(i, j int) bool {
	if res[i].Model < res[j].Model {
		return true
	}
	if res[i].Model > res[j].Model {
		return false
	}
	if res[i].Version < res[j].Version {
		return true
	}
	if res[i].Version > res[j].Version {
		return false
	}
	return res[i].FWCRC < res[j].FWCRC
}
