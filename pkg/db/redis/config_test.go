package redis

import (
	"os"
	"testing"

	"github.com/kelseyhightower/envconfig"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	var config Config

	require.NoError(t, envconfig.Process("", &config))
	assert.Equal(t, config.Hosts, "127.0.0.1")
	assert.Equal(t, config.Port, 6379)
	assert.Equal(t, config.Database, 0)
	assert.Equal(t, config.Username, "")
	assert.Equal(t, config.Password, "")
}

func TestConfig_DSN(t *testing.T) {
	strPtr := func(s string) *string {
		return &s
	}

	type testCase struct {
		Host, Port, Database, Username, Password *string
		ExpectedDSN                              string
	}

	for name, tc := range map[string]testCase{
		"host": {
			Host:        strPtr("0.0.0.0"),
			ExpectedDSN: "redis://0.0.0.0:6379",
		},
		"port": {
			Port:        strPtr("1234"),
			ExpectedDSN: "redis://127.0.0.1:1234",
		},
		"database": {
			Database:    strPtr("3"),
			ExpectedDSN: "redis://127.0.0.1:6379/3",
		},
		"username": {
			Username:    strPtr("user"),
			ExpectedDSN: "redis://user@127.0.0.1:6379",
		},
		"password": {
			Password:    strPtr("secret"),
			ExpectedDSN: "redis://:secret@127.0.0.1:6379",
		},
		"user+password": {
			Username:    strPtr("user"),
			Password:    strPtr("secret"),
			ExpectedDSN: "redis://user:secret@127.0.0.1:6379",
		},
	} {
		t.Run(name, func(t *testing.T) {
			var config Config

			if tc.Host != nil {
				_ = os.Setenv("HOST", *tc.Host)
				defer os.Unsetenv("HOST")
			}
			if tc.Port != nil {
				_ = os.Setenv("PORT", *tc.Port)
				defer os.Unsetenv("PORT")
			}
			if tc.Database != nil {
				_ = os.Setenv("DATABASE", *tc.Database)
				defer os.Unsetenv("DATABASE")
			}
			if tc.Username != nil {
				_ = os.Setenv("USERNAME", *tc.Username)
				defer os.Unsetenv("USERNAME")
			}
			if tc.Password != nil {
				_ = os.Setenv("PASSWORD", *tc.Password)
				defer os.Unsetenv("PASSWORD")
			}

			require.NoError(t, envconfig.Process("", &config))
			assert.Equal(t, config.DSN(), tc.ExpectedDSN)
		})
	}

}
