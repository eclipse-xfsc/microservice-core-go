package goadec_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eclipse-xfsc/microservice-core-go/pkg/goadec"
)

func TestBytesDecoder_Decode(t *testing.T) {
	reqBody := []byte("hello")
	req1, err := http.NewRequest("POST", "https://example.com", bytes.NewReader(reqBody))
	assert.NoError(t, err)

	var v []byte
	err = goadec.BytesDecoder(req1).Decode(&v)
	assert.NoError(t, err)
	assert.Equal(t, reqBody, v)

	req2Body := []byte("hello 2")
	req2, err := http.NewRequest("POST", "https://example.com", bytes.NewReader(req2Body))
	assert.NoError(t, err)

	var vs string
	err = goadec.BytesDecoder(req2).Decode(&vs)
	assert.NoError(t, err)
	assert.Equal(t, string(req2Body), vs)
}
