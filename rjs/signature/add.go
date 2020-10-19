package signature

import "htdvisser.dev/exp/rjs/types"

const AddURI = "/api/v1/signature/add"

type AddRequest struct {
	OwnerID   types.ID6 `json:"ownerid"`
	FWCRC     int32     `json:"fwcrc"`
	KeyCRC    int32     `json:"keycrc"`
	Signature []byte    `json:"signature"`
}

type AddResponse struct {
	SigCRC int32  `json:"sigcrc"`
	Error  string `json:"error,omitempty"`
}
