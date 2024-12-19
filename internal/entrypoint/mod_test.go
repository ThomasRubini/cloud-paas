package entrypoint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExampleFunction(t *testing.T) {
	assert.Equal(t, 3, businessLogic(1, 2))
	assert.Equal(t, 5, businessLogic(1, 4))
}
