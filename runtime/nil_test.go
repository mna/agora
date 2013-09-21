package runtime

import (
	"testing"
)

func TestNilAsBool(t *testing.T) {
	res := Nil.Bool()
	if res != false {
		t.Errorf("Nil as bool : expected %v, got %v", false, res)
	}
}

func TestNilAsString(t *testing.T) {
	res := Nil.String()
	if res != NilString {
		t.Errorf("Nil as string : expected %s, got %s", NilString, res)
	}
}
