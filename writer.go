package collector

type Writer interface {
	WriteFromChannel(channel chan map[string]string)
}
