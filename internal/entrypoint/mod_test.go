package entrypoint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExampleFunction(t *testing.T) {
	// Example test case
	result := businessLogic(1, 2)
	assert.Equal(t, 3, result)
}
