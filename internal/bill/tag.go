package bill

import (
	"regexp"
)

type Tag string

func (t Tag) String() string {
	return string(t)
}

// Validate Tag to be in format 'tag1,tag2,tag3'
func (t Tag) Validate() bool {
	regex := regexp.MustCompile(`([a-z],)`)
	return regex.MatchString(string(t))
}
