package message

type Welcome struct {
	UserName string `json:"user_name"`
	Title    string `json:"title"`
}

func NewWelcome(userName string, title string) *Welcome {
	return &Welcome{
		UserName: userName,
		Title:    title,
	}
}
