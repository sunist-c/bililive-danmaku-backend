package pool

type Pool struct {
	Danmaku  chan []byte
	Gift     chan []byte
	Guard    chan []byte
	Master   chan []byte
	Audience chan []byte
	Unknown  chan []byte
}

func NewPoolWithHandler(handler func(pool *Pool)) *Pool {
	pool := &Pool{
		Danmaku:  make(chan []byte, 32),
		Gift:     make(chan []byte, 32),
		Guard:    make(chan []byte, 32),
		Master:   make(chan []byte, 32),
		Audience: make(chan []byte, 32),
		Unknown:  make(chan []byte, 32),
	}
	go handler(pool)

	return pool
}
