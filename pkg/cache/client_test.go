package cache_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	errors "github.com/eclipse-xfsc/microservice-core-go/pkg/err"
	"github.com/stretchr/testify/assert"

	"github.com/eclipse-xfsc/microservice-core-go/pkg/cache"
)

func TestClient_InvalidCacheAddress(t *testing.T) {
	client := cache.New("dkaslkfasdlkn")

	// get
	res, err := client.Get(context.Background(), "test", "test", "test")
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.True(t, errors.Is(errors.Internal, err))
	assert.Contains(t, err.Error(), "invalid cache url")

	// set
	err = client.Set(context.Background(), "test", "test", "test", []byte("data"))
	assert.Error(t, err)
	assert.True(t, errors.Is(errors.Internal, err))
	assert.Contains(t, err.Error(), "invalid cache url")
}

func TestClient_Get(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		namespace string
		scope     string
		handler   http.HandlerFunc

		result  []byte
		errkind errors.Kind
		errtext string
	}{
		{
			name:      "cache entry not found",
			key:       "mykey",
			namespace: "mynamespace",
			scope:     "myscope",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "mykey", r.Header.Get("x-cache-key"))
				assert.Equal(t, "mynamespace", r.Header.Get("x-cache-namespace"))
				assert.Equal(t, "myscope", r.Header.Get("x-cache-scope"))
				w.WriteHeader(http.StatusNotFound)
			},
			errkind: errors.NotFound,
			errtext: "not found",
		},
		{
			name:      "unexpected error returned from cache",
			key:       "mykey",
			namespace: "mynamespace",
			scope:     "myscope",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "mykey", r.Header.Get("x-cache-key"))
				assert.Equal(t, "mynamespace", r.Header.Get("x-cache-namespace"))
				assert.Equal(t, "myscope", r.Header.Get("x-cache-scope"))
				w.WriteHeader(http.StatusInternalServerError)
			},
			errkind: errors.Internal,
			errtext: "unexpected response: 500 Internal Server Error",
		},
		{
			name:      "cache entry is retrieved successfully",
			key:       "mykey",
			namespace: "mynamespace",
			scope:     "myscope",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "mykey", r.Header.Get("x-cache-key"))
				assert.Equal(t, "mynamespace", r.Header.Get("x-cache-namespace"))
				assert.Equal(t, "myscope", r.Header.Get("x-cache-scope"))
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("data"))
			},
			result: []byte("data"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cachesrv := httptest.NewServer(test.handler)
			client := cache.New(cachesrv.URL)
			result, err := client.Get(context.Background(), test.key, test.namespace, test.scope)
			if test.errtext != "" {
				assert.Nil(t, result)
				assert.Error(t, err)
				assert.True(t, errors.Is(test.errkind, err))
				assert.Contains(t, err.Error(), test.errtext)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.result, result)
			}
		})
	}
}

func TestClient_Set(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		namespace string
		scope     string
		data      []byte
		handler   http.HandlerFunc

		result  []byte
		errkind errors.Kind
		errtext string
	}{
		{
			name:      "unexpected response returned from cache",
			key:       "mykey",
			namespace: "mynamespace",
			scope:     "myscope",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "mykey", r.Header.Get("x-cache-key"))
				assert.Equal(t, "mynamespace", r.Header.Get("x-cache-namespace"))
				assert.Equal(t, "myscope", r.Header.Get("x-cache-scope"))
				w.WriteHeader(http.StatusInternalServerError)
			},
			errkind: errors.Internal,
			errtext: "unexpected response: 500 Internal Server Error",
		},
		{
			name:      "unexpected response returned from cache",
			key:       "mykey",
			namespace: "mynamespace",
			scope:     "myscope",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "mykey", r.Header.Get("x-cache-key"))
				assert.Equal(t, "mynamespace", r.Header.Get("x-cache-namespace"))
				assert.Equal(t, "myscope", r.Header.Get("x-cache-scope"))
				w.WriteHeader(http.StatusRequestTimeout)
			},
			errkind: errors.Timeout,
			errtext: "unexpected response: 408 Request Timeout",
		},
		{
			name:      "data is stored in cache successfully",
			key:       "mykey",
			namespace: "mynamespace",
			scope:     "myscope",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "mykey", r.Header.Get("x-cache-key"))
				assert.Equal(t, "mynamespace", r.Header.Get("x-cache-namespace"))
				assert.Equal(t, "myscope", r.Header.Get("x-cache-scope"))
				w.WriteHeader(http.StatusOK)
			},
		},
		{
			name:      "data is stored in cache successfully",
			key:       "mykey",
			namespace: "mynamespace",
			scope:     "myscope",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "mykey", r.Header.Get("x-cache-key"))
				assert.Equal(t, "mynamespace", r.Header.Get("x-cache-namespace"))
				assert.Equal(t, "myscope", r.Header.Get("x-cache-scope"))
				w.WriteHeader(http.StatusCreated)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cachesrv := httptest.NewServer(test.handler)
			client := cache.New(cachesrv.URL)
			err := client.Set(context.Background(), test.key, test.namespace, test.scope, test.data)
			if test.errtext != "" {
				assert.Error(t, err)
				assert.True(t, errors.Is(test.errkind, err))
				assert.Contains(t, err.Error(), test.errtext)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
