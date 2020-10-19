package firmware

import "htdvisser.dev/exp/rjs/types"

const AddURI = "/api/v1/firmware/add"

type AddRequest struct {
	OwnerID  types.ID6 `json:"ownerid"`
	Firmware []byte    `json:"firmware"`
	Version  string    `json:"version"`
	FWCRC    int32     `json:"fwcrc,omitempty"` // Optional. For adding a MiniHub router firmware, the fwcrc MUST be specified in the firmware/add call. The MiniHub .bin firmware file includes the CRC in the last 4 bytes.
}

type AddResponse struct {
	FWCRC int32  `json:"fwcrc"`
	Error string `json:"error,omitempty"`
}
