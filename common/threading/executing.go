package threading

type Executable interface {
	Execute(exit chan struct{})
	Stop() (success bool)
}
