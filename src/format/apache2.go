package format

type Apache2Config struct {
	Type string
}

const apache2commonRegExp = "^(?P<host>[\\d.]+) (?P<identd>\\S+) (?P<user>\\S+) \\[(?P<time>[\\w:/]+\\s[+\\-]\\d{4})\\] \"(?P<method>\\S+) (?P<path>\\S+) (?P<version>.+?)\" (?P<status>\\d{3}) (?P<size>\\d+)$"
const apache2commonFormat = "(time:\"02/Jan/2006:15:04:05 -0700\")time,(int)status,(int)size"

const apache2combinedRegExp = "^(?P<host>[\\d.]+) (?P<identd>\\S+) (?P<user>\\S+) \\[(?P<time>[\\w:/]+\\s[+\\-]\\d{4})\\] \"(?P<method>\\S+) (?P<path>\\S+) (?P<version>.+?)\" (?P<status>\\d{3}) (?P<size>\\d+) \"(?P<referer>[^\"]+)\" \"(?P<agent>[^\"]+)\"$"
const apache2combinedFormat = "(time:\"02/Jan/2006:15:04:05 -0700\")time,(int)status,(int)size"

const apache2errorRegExp = "^\\[(?P<time>[^\\]]+)\\] \\[(?P<severity>\\S+)\\] \\[(?P<identifier>[^\\]]+)\\] (?P<message>[^\"]+)$"
const apache2errorFormat = "(time:\"Mon Jan 02 15:04:05 2006\")time"

type Apache2 struct {
	RegExp
}

func NewApache2(config *Apache2Config) *Apache2 {
	format := new(Apache2)
	format.SetConfig(format.TransformConfig(config))

	return format
}

func (self *Apache2) TransformConfig(config *Apache2Config) *RegExpConfig {

	var pattern, format string
	switch config.Type {
	case "common":
		pattern = apache2commonRegExp
		format = apache2commonFormat
	case "combined":
		pattern = apache2combinedRegExp
		format = apache2combinedFormat
	case "error":
		pattern = apache2errorRegExp
		format = apache2errorFormat

	}

	regExpConfig := RegExpConfig{Pattern: pattern, Format: format}

	return &regExpConfig
}
