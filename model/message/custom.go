package message

type Custom struct {
	Message string `json:"message"`
}

func NewCustom(message string) *Custom {
	return &Custom{Message: message}
}
