package timer

type Timer interface {
	Start() error
	Stop() error
	Pause() error
	Resume() error
}