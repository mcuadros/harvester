package format

type Apache2Config struct {
	Type string
}

const common = "^(?P<host>[\\d.]+) (?P<identd>\\S+) (?P<user>\\S+) \\[(?P<time>[\\w:/]+\\s[+\\-]\\d{4})\\] \"(?P<method>\\S+) (?P<path>\\S+) (?P<version>.+?)\" (?P<status>\\d{3}) (?P<size>\\d+)$"
const combined = "^(?P<host>[\\d.]+) (?P<identd>\\S+) (?P<user>\\S+) \\[(?P<time>[\\w:/]+\\s[+\\-]\\d{4})\\] \"(?P<method>\\S+) (?P<path>\\S+) (?P<version>.+?)\" (?P<status>\\d{3}) (?P<size>\\d+) \"(?P<referer>[^\"]+)\" \"(?P<agent>[^\"]+)\"$"

type Apache2 struct {
	RegExp
}

func NewApache2(config *Apache2Config) *Apache2 {
	format := new(Apache2)
	format.SetConfig(format.TransformConfig(config))

	return format
}

func (self *Apache2) TransformConfig(config *Apache2Config) *RegExpConfig {

	var pattern string
	switch config.Type {
	case "common":
		pattern = common
	case "combined":
		pattern = combined
	}

	regExpConfig := RegExpConfig{Pattern: pattern}

	return &regExpConfig
}
