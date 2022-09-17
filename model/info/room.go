package info

type Room struct {
	RoomID     uint32 `json:"room_id"`
	UpUID      uint32 `json:"up_uid"`
	Title      string `json:"title"`
	Online     uint32 `json:"online"`
	Tags       string `json:"tags"`
	LiveStatus bool   `json:"live_status"`
	LockStatus bool   `json:"lock_status"`
}
