package router

import (
	"htdvisser.dev/exp/rjs/types"
)

const InfoURI = "/api/v1/router/info"

type InfoRequest struct {
	OwnerID types.ID6 `json:"ownerid"`
	Router  types.ID6 `json:"router"`
}

type InfoBulkRequest struct {
	OwnerID types.ID6   `json:"ownerid"`
	Routers []types.ID6 `json:"routers"` // maximum 128 entries
}

type InfoResponseItem struct {
	Router      types.ID6   `json:"router"`
	Error       string      `json:"error,omitempty"`
	OwnerID     types.ID6   `json:"ownerid,omitempty"`
	CUPSURI     string      `json:"cupsUri,omitempty"`
	CUPSCredCRC int32       `json:"cupsCredCrc,omitempty"`
	CUPSKey     []byte      `json:"cupsKey,omitempty"`
	CUPSCrt     []byte      `json:"cupsCrt,omitempty"`
	CUPSTrust   []byte      `json:"cupsTrust,omitempty"`
	LNSURI      string      `json:"lnsUri,omitempty"`
	LNSCredCRC  int32       `json:"lnsCredCrc,omitempty"`
	LNSKey      []byte      `json:"lnsKey,omitempty"`
	LNSCrt      []byte      `json:"lnsCrt,omitempty"`
	LNSTrust    []byte      `json:"lnsTrust,omitempty"`
	Claimed     *types.Time `json:"claimed,omitempty"`
	LastSetup   *types.Time `json:"lastSetup,omitempty"`
	LastSeen    *types.Time `json:"lastSeen,omitempty"`
	HTTPStatus  uint16      `json:"httpStatus,omitempty"`
}

type InfoResponse []*InfoResponseItem
