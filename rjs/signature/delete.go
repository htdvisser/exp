package signature

import "htdvisser.dev/exp/rjs/types"

const DeleteURI = "/api/v1/signature/delete"

type DeleteRequest struct {
	OwnerID types.ID6 `json:"ownerid"`
	SigCRC  int32     `json:"sigcrc"`
}

type DeleteResponse struct {
	Error string `json:"error,omitempty"`
}
