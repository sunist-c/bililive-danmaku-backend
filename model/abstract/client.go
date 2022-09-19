package abstract

type IClient interface {
	Close() (success bool)
	Serve() (success bool)
}
