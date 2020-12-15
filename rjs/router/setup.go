package router

import "htdvisser.dev/exp/rjs/types"

const SetupURI = "/api/v1/router/setup"

type SetupRequestItem struct {
	Router    types.ID6   `json:"router,omitempty"`
	CUPSURI   string      `json:"cupsUri,omitempty"`
	CUPSKey   []byte      `json:"cupsKey,omitempty"`
	CUPSCrt   []byte      `json:"cupsCrt,omitempty"`
	CUPSTrust []byte      `json:"cupsTrust,omitempty"`
	LNSURI    string      `json:"lnsUri,omitempty"`
	LNSKey    []byte      `json:"lnsKey,omitempty"`
	LNSCrt    []byte      `json:"lnsCrt,omitempty"`
	LNSTrust  []byte      `json:"lnsTrust,omitempty"`
	FWCRC     int32       `json:"fwcrc,omitempty"`
	FWAfter   *types.Time `json:"fwafter,omitempty"`
}

type SetupRequest struct {
	OwnerID types.ID6 `json:"ownerid"`
	SetupRequestItem
}

type SetupBulkRequest struct {
	OwnerID types.ID6          `json:"ownerid"`
	Routers []SetupRequestItem `json:"routers"`
}

type SetupMixedRequest struct {
	OwnerID types.ID6 `json:"ownerid"`
	SetupRequestItem
	Routers []SetupRequestItem `json:"routers"`
}

type SetupResponseItem struct {
	Router types.ID6 `json:"router"`
	Error  string    `json:"error,omitempty"`
}

type SetupResponse []*SetupResponseItem
