package tests

import (
	"testing"
)

func TestExample(t *testing.T) {
	if 1+1 != 2 {
		t.Error("Expected 1 + 1 to equal 2")
	}
}
