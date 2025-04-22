package ocm_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eclipse-xfsc/microservice-core-go/pkg/ocm"
)

func Test_GetLoginProofInvitationSuccess(t *testing.T) {
	expected := &ocm.LoginProofInvitationResponse{
		StatusCode: 200,
		Message:    "success",
		Data:       ocm.LoginProofInvitationResponseData{},
	}

	ocmServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(expected)
	}))

	client := ocm.New(ocmServer.URL)
	res, err := client.GetLoginProofInvitation(context.Background(), []string{"principalMembershipCredential"})

	assert.NoError(t, err)
	assert.Equal(t, expected.StatusCode, res.StatusCode)
	assert.Equal(t, expected.Message, res.Message)
	assert.Equal(t, expected.Data, res.Data)
}

func Test_GetLoginProofInvitationErr(t *testing.T) {
	ocmServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	client := ocm.New(ocmServer.URL)
	res, err := client.GetLoginProofInvitation(context.Background(), []string{"principalMembershipCredential"})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected response code")
}

func Test_SendOutOfBandRequestSuccess(t *testing.T) {
	expected := &ocm.LoginProofInvitationResponse{
		StatusCode: 200,
		Message:    "success",
		Data:       ocm.LoginProofInvitationResponseData{},
	}

	ocmServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(expected)
	}))

	client := ocm.New(ocmServer.URL)
	res, err := client.SendOutOfBandRequest(context.Background(), map[string]interface{}{
		"attributes": []map[string]string{
			{
				"schemaId":               "7KuDTpQh3GJ7Gp6kErpWvM:2:principalTestSchema:1.0",
				"credentialDefinitionId": "7KuDTpQh3GJ7Gp6kErpWvM:3:CL:40329:principalTestCredDefExpir",
				"attributeName":          "prcLastName",
				"value":                  "",
			},
			{
				"schemaId":               "7KuDTpQh3GJ7Gp6kErpWvM:2:principalTestSchema:1.0",
				"credentialDefinitionId": "7KuDTpQh3GJ7Gp6kErpWvM:3:CL:40329:principalTestCredDefExpir",
				"attributeName":          "email",
				"value":                  "",
			},
		},
		"options": map[string]string{
			"type": "Aries1.0",
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, expected.StatusCode, res.StatusCode)
	assert.Equal(t, expected.Message, res.Message)
	assert.Equal(t, expected.Data, res.Data)
}

func Test_SendOutOfBandRequestErr(t *testing.T) {
	ocmServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	client := ocm.New(ocmServer.URL)
	res, err := client.SendOutOfBandRequest(context.Background(), map[string]interface{}{
		"attributes": []map[string]string{
			{
				"schemaId":               "7KuDTpQh3GJ7Gp6kErpWvM:2:principalTestSchema:1.0",
				"credentialDefinitionId": "7KuDTpQh3GJ7Gp6kErpWvM:3:CL:40329:principalTestCredDefExpir",
				"attributeName":          "prcLastName",
				"value":                  "",
			},
			{
				"schemaId":               "7KuDTpQh3GJ7Gp6kErpWvM:2:principalTestSchema:1.0",
				"credentialDefinitionId": "7KuDTpQh3GJ7Gp6kErpWvM:3:CL:40329:principalTestCredDefExpir",
				"attributeName":          "email",
				"value":                  "",
			},
		},
		"options": map[string]string{
			"type": "Aries1.0",
		},
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected response code")
}

func TestClient_GetLoginProofResultSuccess(t *testing.T) {
	expected := &ocm.LoginProofResultResponse{
		StatusCode: 200,
		Message:    "success",
		Data:       ocm.LoginProofResultResponseData{},
	}

	ocmServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(expected)
	}))

	client := ocm.New(ocmServer.URL)
	res, err := client.GetLoginProofResult(context.Background(), "2cf01406-b15f-4960-a6a7-7bc62cd37a3c")

	assert.NoError(t, err)
	assert.Equal(t, expected.StatusCode, res.StatusCode)
	assert.Equal(t, expected.Message, res.Message)
	assert.Equal(t, expected.Data, res.Data)
}

func Test_GetLoginProofResultErr(t *testing.T) {
	ocmServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	client := ocm.New(ocmServer.URL)
	res, err := client.GetLoginProofResult(context.Background(), "2cf01406-b15f-4960-a6a7-7bc62cd37a3c")

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected response code")
}

func Test_GetRawLoginProofResultSuccess(t *testing.T) {
	expected := map[string]interface{}{
		"statusCode": float64(200),
		"message":    "Proof presentation fetch successfully",
		"data": map[string]interface{}{
			"state": "done",
			"presentations": []interface{}{
				map[string]interface{}{
					"schemaId":  "7KuDTpQh3GJ7Gp6kErpWvM:2:principalTestSchema:1.0",
					"credDefId": "7KuDTpQh3GJ7Gp6kErpWvM:3:CL:40329:principalTestCredDefExpire",
					"revRegId":  nil,
					"timestamp": nil,
					"credentialSubject": map[string]interface{}{
						"email":       "23957edb-991d-4b5f-bf76-153103ba45b7",
						"prcLastName": "NA",
					},
				},
			},
		},
	}

	ocmServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(expected)
	}))

	client := ocm.New(ocmServer.URL)
	res, err := client.GetRawLoginProofResult(context.Background(), "2cf01406-b15f-4960-a6a7-7bc62cd37a3c")

	assert.NoError(t, err)
	assert.Equal(t, expected["statusCode"], res["statusCode"])
	assert.Equal(t, expected["message"], res["message"])
	assert.Equal(t, expected["data"], res["data"])
}

func Test_GetRawLoginProofResultErr(t *testing.T) {
	ocmServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	client := ocm.New(ocmServer.URL)
	res, err := client.GetRawLoginProofResult(context.Background(), "2cf01406-b15f-4960-a6a7-7bc62cd37a3c")

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected response code")
}
