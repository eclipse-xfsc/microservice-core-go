package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eclipse-xfsc/microservice-core-go/pkg/auth"
)

var (
	publicKey  jwk.RSAPublicKey
	privateKey jwk.RSAPrivateKey
)

// initKeys creates private and public RSA keys and sets them
// in global variables that are used by all tests.
func initKeys() error {
	rawprivkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to create raw private key: %v", err)
	}

	privkey, err := jwk.FromRaw(rawprivkey)
	if err != nil {
		return fmt.Errorf("failed to create private key: %v", err)
	}

	pubkey, err := privkey.PublicKey()
	if err != nil {
		return fmt.Errorf("failed to create public key: %v", err)
	}

	privk, ok := privkey.(jwk.RSAPrivateKey)
	if !ok {
		return fmt.Errorf("cannot cast private key to RSA private key")
	}
	privateKey = privk

	if err := privateKey.Set(jwk.KeyIDKey, "key1"); err != nil {
		return fmt.Errorf("cannot set kid value to private key: %v", err)
	}

	pubk, ok := pubkey.(jwk.RSAPublicKey)
	if !ok {
		return fmt.Errorf("cannot cast public key to RSA public key")
	}
	publicKey = pubk

	return nil
}

func TestAuthMiddleware_Handler(t *testing.T) {
	err := initKeys()
	require.NoError(t, err)

	keyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a JWK Set
		set := jwk.NewSet()

		var raw interface{}
		err := publicKey.Raw(&raw)
		assert.NoError(t, err)

		key, err := jwk.FromRaw(raw)
		assert.NoError(t, err)

		err = key.Set(jwk.AlgorithmKey, jwa.RS256)
		assert.NoError(t, err)

		err = key.Set("kid", "key1")
		assert.NoError(t, err)

		err = set.AddKey(key)
		assert.NoError(t, err)

		err = json.NewEncoder(w).Encode(set)
		assert.NoError(t, err)
	}))

	authMiddleware, err := auth.NewMiddleware(keyServer.URL, 1*time.Hour, http.DefaultClient)
	require.NoError(t, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("everything is fine"))
	})
	authHandler := authMiddleware.Handler()(handler)

	t.Run("authenticate with valid token", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
		assert.NoError(t, err)

		token, err := createSignedToken()
		assert.NoError(t, err)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		response := httptest.NewRecorder()
		authHandler.ServeHTTP(response, req)

		assert.Equal(t, "everything is fine", response.Body.String())
	})

	t.Run("authenticate with invalid token", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
		assert.NoError(t, err)

		token := "deadbeef"
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		response := httptest.NewRecorder()
		authHandler.ServeHTTP(response, req)

		assert.Equal(t, "failed to parse jws: invalid compact serialization format: invalid number of segments\n", response.Body.String())
	})

	t.Run("authenticate with token signed with unknown key (invalid signature)", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
		assert.NoError(t, err)
		// gg-ignore: test secret, not valid
		token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6ImtleTEifQ.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTY2NDk2MDg2OCwiZXhwIjoxNjY0OTY0NDY4fQ.FrIA3A228den86qX72o4yP3TiEA9uOf46Yav_vY5daCQ8yeAm3GaBzC_nikt0y9NSCR6K2G2GCm7RcdfP3vQ9CFh2R7FtL4nfjffdauLmXVzp3z_lyBIKYL3RsTGChctfMeYZzk2F6EDmGHeI8xV3KiDC5Gfkvfdp9MfFxVy7DcuEV9MLo_9j4Y-7nfuB1CbdF_1vzSsO0twitePjsB59CNndugJgTUGFjKUJU2_e7vKMR_i9NvFHfJZS2VbtX3vrZ5f_pfOvBSSZJBxG50Uwf6COhtABieVHhhmLBSJq1P1EWRAI26Bk-YtE8k-jfjra9W1RF5DLF7Jh9Lw-utc5A" //nolint:gosec
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		response := httptest.NewRecorder()
		authHandler.ServeHTTP(response, req)

		assert.Equal(t, "could not verify message using any of the signatures or keys\n", response.Body.String())
	})

	t.Run("request without token", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
		assert.NoError(t, err)

		response := httptest.NewRecorder()
		authHandler.ServeHTTP(response, req)

		assert.Equal(t, "invalid authorization header\n", response.Body.String())
	})

	t.Run("invalid authorization header", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
		assert.NoError(t, err)

		token, err := createSignedToken()
		assert.NoError(t, err)
		req.Header.Set("Authorization", token)

		response := httptest.NewRecorder()
		authHandler.ServeHTTP(response, req)

		assert.Equal(t, "invalid authorization header\n", response.Body.String())
	})
}

func createSignedToken() (string, error) {
	token, err := jwt.NewBuilder().
		Claim(`claim1`, `value1`).
		Claim(`claim2`, `value2`).
		Issuer(`https://example.com`).
		Subject("terminator").
		Audience([]string{"skynet"}).
		Build()
	if err != nil {
		return "", fmt.Errorf("failed to build token: %s", err)
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, privateKey))
	if err != nil {
		return "", err
	}

	return string(signed), nil
}
