package format

type NginxConfig struct {
	Type string
}

const nginxErrorRegExp = "^(?P<time>[\\d+/ :]+) \\[(?P<severity>.+)\\] .*?: (?P<message>.+), client: (?P<client>.+), server: (?P<server>.+), request: \"(?P<method>\\S+) (?P<path>\\S+) (?P<version>.+?)\", host: \"(?P<host>.+)\"$"
const nginxErrorFormat = "(time:\"2006/01/02 15:04:05\")time"

type Nginx struct {
	RegExp
}

func NewNginx(config *NginxConfig) *Nginx {
	format := new(Nginx)
	format.SetConfig(format.TransformConfig(config))

	return format
}

func (self *Nginx) TransformConfig(config *NginxConfig) *RegExpConfig {

	var pattern, format string
	switch config.Type {
	case "combined":
		pattern = apache2combinedRegExp
		format = apache2commonFormat

	case "error":
		pattern = nginxErrorRegExp
		format = nginxErrorFormat
	}

	regExpConfig := RegExpConfig{Pattern: pattern, Format: format}

	return &regExpConfig
}
