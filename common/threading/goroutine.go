package threading

import (
	"github.com/sunist-c/bililive-danmaku-backend/common/logging"
	"runtime"
)

type Goroutine struct {
	exitChan chan struct{}
	function Executable
}

func (g *Goroutine) Close() (success bool) {
	pc, _, _, _ := runtime.Caller(2)
	class := runtime.FuncForPC(pc).Name()
	logging.Info("exit current goroutine: %+v by %s", g.function, class)
	g.exitChan <- struct{}{}
	return g.function.Stop()
}

func (g *Goroutine) Serve() (success bool) {
	go g.function.Execute(g.exitChan)
	return true
}

func NewGoroutine(function Executable) *Goroutine {
	return &Goroutine{
		exitChan: make(chan struct{}),
		function: function,
	}
}
