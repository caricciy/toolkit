package toolkit

import "testing"

func TestTools_RandomString(t *testing.T) {
	var (
		testTools Tools
	)

	s := testTools.RandomString(10)

	if len(s) != 10 {
		t.Errorf("RandomString() = %v, want %v", len(s), 10)
	}
}