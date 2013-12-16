package collector

type Format interface {
	Parse(line string) map[string]string
}
