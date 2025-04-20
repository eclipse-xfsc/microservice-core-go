package environment

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewDefaultEnv(t *testing.T) {
	env := NewDefaultEnv()
	require.NotNil(t, env)
	assert.IsType(t, &DefaultEnv{}, env)
}

func TestDefaultEnv_IsHealthy(t *testing.T) {
	env := NewDefaultEnv()
	assert.True(t, env.IsHealthy())
}

func TestDefaultEnv_SetHealthFunc(t *testing.T) {
	env := NewDefaultEnv()
	called := false

	hf := func() bool {
		called = true
		return true
	}

	env.SetHealthFunc(hf)

	assert.True(t, env.IsHealthy())
	assert.True(t, called, "the provided health func ")
}
