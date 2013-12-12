package collector

type Reader interface {
	ReadIntoChannel(channel chan map[string]string)
}
