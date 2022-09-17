package model

type Api string

func (u Api) ToString() string {
	return string(u)
}
