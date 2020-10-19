package firmware

import "htdvisser.dev/exp/rjs/types"

const DeleteURI = "/api/v1/firmware/delete"

type DeleteRequest struct {
	OwnerID types.ID6 `json:"ownerid"`
	FWCRC   int32     `json:"fwcrc"`
}

type DeleteResponse struct {
	Error string `json:"error,omitempty"`
}
