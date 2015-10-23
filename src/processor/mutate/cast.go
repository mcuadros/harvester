package mutate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	INT  = "int"
	DATE = "date"
)

func cast(value interface{}, fn string, params ...string) (interface{}, error) {
	switch fn {
	case INT:
		switch value.(type) {
		case string:
			for _, param := range params {
				switch param {
				case "strip":
					re := regexp.MustCompile("[0-9]+")
					value = re.FindString(value.(string))
				}
			}
			result, err := strconv.Atoi(value.(string))
			if err != nil {
				return value, err
			} else {
				return result, nil
			}
		case int:
		default:
			return value, fmt.Errorf("cast: invalid input value `%s`", value)
		}
	case DATE:
		switch value.(type) {
		case string:
			if len(params) < 1 {
				return value, fmt.Errorf("cast: cast string to date requires a format, can't cast `%s`", value)
			}
			var err error
			var result time.Time
			for _, format := range params {
				if format == "null" {
					return nil, nil
				}
				if format == "present" && strings.ToLower(value.(string)) == "present" {
					return time.Now(), nil
				}
				result, err = time.Parse(format, value.(string))
				if err == nil {
					break
				}
			}
			if err != nil {
				return value, err
			} else {
				return result, nil
			}
		case int:
			return time.Unix(int64(value.(int)), 0), nil
		default:
			return value, fmt.Errorf("cast: invalid input value `%s`", value)
		}
	default:
		return value, fmt.Errorf("cast: unrecognized cast function `%s`", fn)
	}

	return value, nil
}
