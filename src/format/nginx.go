package format

type NginxConfig struct {
	Type string
}

const nginxError = "^(?P<time>[\\d+/ :]+) \\[(?P<severity>.+)\\] .*?: (?P<message>.+), client: (?P<client>.+), server: (?P<server>.+), request: \"(?P<method>\\S+) (?P<path>\\S+) (?P<version>.+?)\", host: \"(?P<host>.+)\"$"

type Nginx struct {
	RegExp
}

func NewNginx(config *NginxConfig) *Nginx {
	format := new(Nginx)
	format.SetConfig(format.TransformConfig(config))

	return format
}

func (self *Nginx) TransformConfig(config *NginxConfig) *RegExpConfig {

	var pattern string
	switch config.Type {
	case "combined":
		pattern = apache2combined
	case "error":
		pattern = nginxError
	}

	regExpConfig := RegExpConfig{Pattern: pattern}

	return &regExpConfig
}
