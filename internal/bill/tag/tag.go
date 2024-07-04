package tag

import (
	"regexp"
)

type Tag struct {
	String string
	Valid  bool
}

func New(tag string) *Tag {
	if tag == "" {
		return &Tag{
			String: "empty",
			Valid:  false,
		}
	}
	return &Tag{
		String: tag,
		Valid:  true,
	}
}

func NewFromNullable(tag *string) *Tag {
	if tag == nil {
		return &Tag{
			String: "empty",
			Valid:  false,
		}
	}
	return &Tag{
		String: *tag,
		Valid:  true,
	}
}

func Empty() *Tag {
	return &Tag{
		String: "empty",
		Valid:  false,
	}
}

// Validate Tag to be in format 'tag1,tag2,tag3'
func (t Tag) Validate() bool {
	regex := regexp.MustCompile(`([a-z],)`)
	return regex.MatchString(t.String)
}
