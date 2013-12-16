package collector

type Input interface {
	GetLine() string
	IsEOF() bool
}
