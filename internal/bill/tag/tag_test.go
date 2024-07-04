package tag

import "testing"

func TestTagNew(t *testing.T) {
	tagString := "tag1,tag2,tag3"
	tag := New(tagString)
	if tag.String != "tag1,tag2,tag3" {
		t.Errorf("Tag string is not correct: %s\n", tag.String)
	}
	if !tag.Valid {
		t.Errorf("Tag should be valid: %s\n", tagString)
	}
	tag = New("")
	if tag.Valid {
		t.Error("Tag should be not valid")
	}
}

func TestTagNewFromNullable(t *testing.T) {
	tagString := "tag1,tag2,tag3"
	tag := NewFromNullable(&tagString)
	if tag.String != "tag1,tag2,tag3" {
		t.Errorf("Tag string is not correct: %s\n", tag.String)
	}
	if !tag.Valid {
		t.Errorf("Tag should be valid: %s\n", tagString)
	}
	tagString = ""
	tag = NewFromNullable(&tagString)
	if tag.Valid {
		t.Errorf(
			"Tag should be not valid. Created from: '%s'\n",
			tagString,
		)
	}

	tag = NewFromNullable(nil)
	if tag.Valid {
		t.Error("Tag should be not valid")
	}
}
