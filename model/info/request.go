package info

type Request struct {
	UID           uint8  `json:"uid"`
	RoomID        uint32 `json:"roomid"`
	ProtoVersion  uint8  `json:"protover"`
	Platform      string `json:"platform"`
	ClientVersion string `json:"clientver"`
	Type          uint8  `json:"type"`
	Key           string `json:"key"`
}
