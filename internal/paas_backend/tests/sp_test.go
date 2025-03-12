package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEmptySecret(t *testing.T) {
	state := fakeState(t)
	value, err := state.SecretsProvider.GetSecret("key1")
	assert.Nil(t, err)
	assert.Equal(t, "", value)
}

func TestSetAndGetSecrets(t *testing.T) {
	state := fakeState(t)

	err := state.SecretsProvider.SetSecret("key1", "value1")
	assert.Nil(t, err)
	err = state.SecretsProvider.SetSecret("key2", "value2")
	assert.Nil(t, err)

	value, err := state.SecretsProvider.GetSecret("key1")
	assert.Nil(t, err)
	assert.Equal(t, "value1", value)

	value, err = state.SecretsProvider.GetSecret("key2")
	assert.Nil(t, err)
	assert.Equal(t, "value2", value)
}

func TestSetSecretTwice(t *testing.T) {
	state := fakeState(t)

	err := state.SecretsProvider.SetSecret("key1", "value1")
	assert.Nil(t, err)

	err = state.SecretsProvider.SetSecret("key1", "value2")
	assert.Nil(t, err)
}
