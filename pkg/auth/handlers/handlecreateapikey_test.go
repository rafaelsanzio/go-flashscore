package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/rafaelsanzio/go-flashscore/pkg/apikey"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
)

func mockCreateApiKeyFunc(ctx context.Context, a apikey.ApiKey) errs.AppError {
	return nil
}

func mockCreateApiKeyThrowFunc(ctx context.Context, a apikey.ApiKey) errs.AppError {
	return errs.ErrRepoMockAction
}

func mockFindByEmailApiKeyFunc(ctx context.Context, email string) (*apikey.ApiKey, errs.AppError) {
	apikey := prototype.PrototypeApiKey()
	return &apikey, nil
}

func mockFindByEmailApiKeyEmptyFunc(ctx context.Context, email string) (*apikey.ApiKey, errs.AppError) {
	return nil, nil
}

func mockFindByEmailApiKeyThrowFunc(ctx context.Context, email string) (*apikey.ApiKey, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func mockUpdateApiKeyFunc(ctx context.Context, a apikey.ApiKey) (*apikey.ApiKey, errs.AppError) {
	apikey := prototype.PrototypeApiKey()
	return &apikey, nil
}

func mockUpdateApiKeyThrowFunc(ctx context.Context, a apikey.ApiKey) (*apikey.ApiKey, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func TestHandleCreateApiKey(t *testing.T) {
	body, err := json.Marshal(APIKeyPayload{
		Email: "rafael@gmail.com",
	})
	assert.Equal(t, nil, err)

	bodyWrongMail, err := json.Marshal(APIKeyPayload{
		Email: "rafaelgmail.com",
	})
	assert.Equal(t, nil, err)

	goodReq := httptest.NewRequest(http.MethodPost, "/apikey", nil)
	goodReq = mux.SetURLVars(goodReq, map[string]string{})
	goodReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	noBodyReq := httptest.NewRequest(http.MethodPost, "/apikey", nil)
	noBodyReq = mux.SetURLVars(noBodyReq, map[string]string{})

	throwReq := httptest.NewRequest(http.MethodPost, "/apikey", nil)
	throwReq = mux.SetURLVars(throwReq, map[string]string{})
	throwReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	goodReq2 := httptest.NewRequest(http.MethodPost, "/apikey", nil)
	goodReq2 = mux.SetURLVars(goodReq2, map[string]string{})
	goodReq2.Body = ioutil.NopCloser(bytes.NewReader(body))

	goodReq3 := httptest.NewRequest(http.MethodPost, "/apikey", nil)
	goodReq3 = mux.SetURLVars(goodReq3, map[string]string{})
	goodReq3.Body = ioutil.NopCloser(bytes.NewReader(body))

	goodReq4 := httptest.NewRequest(http.MethodPost, "/apikey", nil)
	goodReq4 = mux.SetURLVars(goodReq4, map[string]string{})
	goodReq4.Body = ioutil.NopCloser(bytes.NewReader(body))

	goodReq5 := httptest.NewRequest(http.MethodPost, "/apikey", nil)
	goodReq5 = mux.SetURLVars(goodReq5, map[string]string{})
	goodReq5.Body = ioutil.NopCloser(bytes.NewReader(body))

	goodReq6 := httptest.NewRequest(http.MethodPost, "/apikey", nil)
	goodReq6 = mux.SetURLVars(goodReq6, map[string]string{})
	goodReq6.Body = ioutil.NopCloser(bytes.NewReader(body))

	goodReq7 := httptest.NewRequest(http.MethodPost, "/apikey", nil)
	goodReq7 = mux.SetURLVars(goodReq7, map[string]string{})
	goodReq7.Body = ioutil.NopCloser(bytes.NewReader(body))

	goodReq8 := httptest.NewRequest(http.MethodPost, "/apikey", nil)
	goodReq8 = mux.SetURLVars(goodReq8, map[string]string{})
	goodReq8.Body = ioutil.NopCloser(bytes.NewReader(bodyWrongMail))

	goodReq9 := httptest.NewRequest(http.MethodPost, "/apikey", nil)
	goodReq9 = mux.SetURLVars(goodReq9, map[string]string{})
	goodReq9.Body = ioutil.NopCloser(bytes.NewReader(body))

	testCases := []struct {
		Name                       string
		Request                    *http.Request
		HandleCreateFunc           func(ctx context.Context, a apikey.ApiKey) errs.AppError
		HandleFindByEmailFunc      func(ctx context.Context, email string) (*apikey.ApiKey, errs.AppError)
		HandleUpdateFunc           func(ctx context.Context, a apikey.ApiKey) (*apikey.ApiKey, errs.AppError)
		ConvertingPayloadFunc      func(a APIKeyPayload) (apikey.ApiKey, errs.AppError)
		ConfigGenerateJWTTokenFunc func(email string) (string, errs.AppError)
		MarshalFunc                func(v interface{}) ([]byte, error)
		WriteFunc                  func(http.ResponseWriter, []byte) (int, error)
		ExpectedStatusCode         int
	}{
		{
			Name:                       "Should return 200 if successful",
			Request:                    goodReq,
			HandleCreateFunc:           mockCreateApiKeyFunc,
			HandleFindByEmailFunc:      mockFindByEmailApiKeyEmptyFunc,
			HandleUpdateFunc:           nil,
			ConvertingPayloadFunc:      convertPayloadToApiKeyFunc,
			ConfigGenerateJWTTokenFunc: configGenerateJWTToken,
			MarshalFunc:                jsonMarshal,
			WriteFunc:                  write,
			ExpectedStatusCode:         200,
		}, {
			Name:                       "Should return 422 bad request",
			Request:                    noBodyReq,
			HandleCreateFunc:           mockCreateApiKeyFunc,
			HandleFindByEmailFunc:      mockFindByEmailApiKeyEmptyFunc,
			HandleUpdateFunc:           nil,
			ConvertingPayloadFunc:      convertPayloadToApiKeyFunc,
			ConfigGenerateJWTTokenFunc: configGenerateJWTToken,
			MarshalFunc:                jsonMarshal,
			WriteFunc:                  write,
			ExpectedStatusCode:         422,
		}, {
			Name:                       "Should return 422 throwing error on converting payload func",
			Request:                    goodReq2,
			HandleCreateFunc:           mockCreateApiKeyFunc,
			HandleFindByEmailFunc:      mockFindByEmailApiKeyEmptyFunc,
			HandleUpdateFunc:           nil,
			ConvertingPayloadFunc:      fakeConvertPayloadToApiKeyFunc,
			ConfigGenerateJWTTokenFunc: configGenerateJWTToken,
			MarshalFunc:                jsonMarshal,
			WriteFunc:                  write,
			ExpectedStatusCode:         422,
		}, {
			Name:                       "Should return 500 throwing error on create function",
			Request:                    throwReq,
			HandleCreateFunc:           mockCreateApiKeyThrowFunc,
			HandleFindByEmailFunc:      mockFindByEmailApiKeyEmptyFunc,
			HandleUpdateFunc:           nil,
			ConvertingPayloadFunc:      convertPayloadToApiKeyFunc,
			ConfigGenerateJWTTokenFunc: configGenerateJWTToken,
			MarshalFunc:                jsonMarshal,
			WriteFunc:                  write,
			ExpectedStatusCode:         500,
		}, {
			Name:                       "Should return 500 throwing error on marshal function",
			Request:                    goodReq3,
			HandleCreateFunc:           mockCreateApiKeyFunc,
			HandleFindByEmailFunc:      mockFindByEmailApiKeyEmptyFunc,
			HandleUpdateFunc:           nil,
			ConvertingPayloadFunc:      convertPayloadToApiKeyFunc,
			ConfigGenerateJWTTokenFunc: configGenerateJWTToken,
			MarshalFunc:                fakeMarshal,
			WriteFunc:                  write,
			ExpectedStatusCode:         500,
		}, {
			Name:                       "Should return 500 throwing error on write function",
			Request:                    goodReq4,
			HandleCreateFunc:           mockCreateApiKeyFunc,
			HandleFindByEmailFunc:      mockFindByEmailApiKeyEmptyFunc,
			HandleUpdateFunc:           nil,
			ConvertingPayloadFunc:      convertPayloadToApiKeyFunc,
			ConfigGenerateJWTTokenFunc: configGenerateJWTToken,
			MarshalFunc:                jsonMarshal,
			WriteFunc:                  fakeWrite,
			ExpectedStatusCode:         500,
		}, {
			Name:                       "Should return 500 throwing error on find by email function",
			Request:                    goodReq5,
			HandleCreateFunc:           mockCreateApiKeyFunc,
			HandleFindByEmailFunc:      mockFindByEmailApiKeyThrowFunc,
			HandleUpdateFunc:           nil,
			ConvertingPayloadFunc:      convertPayloadToApiKeyFunc,
			ConfigGenerateJWTTokenFunc: configGenerateJWTToken,
			MarshalFunc:                jsonMarshal,
			WriteFunc:                  write,
			ExpectedStatusCode:         500,
		}, {
			Name:                       "Should return 200 if successful",
			Request:                    goodReq6,
			HandleCreateFunc:           mockCreateApiKeyFunc,
			HandleFindByEmailFunc:      mockFindByEmailApiKeyFunc,
			HandleUpdateFunc:           mockUpdateApiKeyFunc,
			ConvertingPayloadFunc:      convertPayloadToApiKeyFunc,
			ConfigGenerateJWTTokenFunc: configGenerateJWTToken,
			MarshalFunc:                jsonMarshal,
			WriteFunc:                  write,
			ExpectedStatusCode:         200,
		}, {
			Name:                       "Should return 500 throwing error on update function",
			Request:                    goodReq7,
			HandleCreateFunc:           mockCreateApiKeyFunc,
			HandleFindByEmailFunc:      mockFindByEmailApiKeyFunc,
			HandleUpdateFunc:           mockUpdateApiKeyThrowFunc,
			ConvertingPayloadFunc:      convertPayloadToApiKeyFunc,
			ConfigGenerateJWTTokenFunc: configGenerateJWTToken,
			MarshalFunc:                jsonMarshal,
			WriteFunc:                  write,
			ExpectedStatusCode:         500,
		}, {
			Name:                       "Should return 422 on wrong email type",
			Request:                    goodReq8,
			HandleCreateFunc:           mockCreateApiKeyFunc,
			HandleFindByEmailFunc:      mockFindByEmailApiKeyEmptyFunc,
			HandleUpdateFunc:           nil,
			ConvertingPayloadFunc:      convertPayloadToApiKeyFunc,
			ConfigGenerateJWTTokenFunc: configGenerateJWTToken,
			MarshalFunc:                jsonMarshal,
			WriteFunc:                  write,
			ExpectedStatusCode:         422,
		}, {
			Name:                       "Should return 422 throwing error on generate JWT",
			Request:                    goodReq9,
			HandleCreateFunc:           mockCreateApiKeyFunc,
			HandleFindByEmailFunc:      mockFindByEmailApiKeyEmptyFunc,
			HandleUpdateFunc:           nil,
			ConvertingPayloadFunc:      convertPayloadToApiKeyFunc,
			ConfigGenerateJWTTokenFunc: fakeConfigGenerateJWTToken,
			MarshalFunc:                jsonMarshal,
			WriteFunc:                  write,
			ExpectedStatusCode:         422,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetApiKeyRepo(repo.MockAPIKeyRepo{
			InsertFunc:      tc.HandleCreateFunc,
			FindByEmailFunc: tc.HandleFindByEmailFunc,
			UpdateFunc:      tc.HandleUpdateFunc,
		})
		defer repo.SetApiKeyRepo(nil)

		convertPayloadToApiKeyFunc = tc.ConvertingPayloadFunc
		defer restoreConvertPayloadToApiKeyFunc(convertPayloadToApiKeyFunc)

		configGenerateJWTToken = tc.ConfigGenerateJWTTokenFunc
		defer restoreConfigGenerateJWTToken(configGenerateJWTToken)

		jsonMarshal = tc.MarshalFunc
		defer restoreMarshal(jsonMarshal)

		write = tc.WriteFunc
		defer restoreWrite(write)

		w := httptest.NewRecorder()

		HandleCreateAPIKey(w, tc.Request)
		t.Logf("Response Body: %v", w.Body)

		assert.Equal(t, tc.ExpectedStatusCode, w.Code)
	}
}

func TestConvertPayloadToApiKey(t *testing.T) {
	inPayload := APIKeyPayload{
		Email: "rafael@gmail.com",
	}

	expectedApiKey := apikey.ApiKey{
		ID:         "",
		Email:      "rafael@gmail.com",
		Key:        "",
		ValidUntil: time.Now().AddDate(0, 0, 10).Truncate(time.Second),
	}

	testCases := []struct {
		Name           string
		Payload        APIKeyPayload
		ExpectedApiKey apikey.ApiKey
		ExpectError    bool
		ExpectedError  string
	}{
		{
			Name:           "Test Case: 1 - correct body, no error",
			Payload:        inPayload,
			ExpectedApiKey: expectedApiKey,
			ExpectError:    false,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		team, err := convertPayloadToApiKey(tc.Payload)
		if tc.ExpectError {
			assert.NotNil(t, err)
			assert.Equal(t, err.Error(), tc.ExpectedError)
		} else {
			assert.Equal(t, tc.ExpectedApiKey, team)
		}
	}
}

func TestDecodeApiKeyRequest(t *testing.T) {
	body, err := json.Marshal(APIKeyPayload{
		Email: "rafael@gmail.com",
	})
	assert.Equal(t, nil, err)

	goodReq := httptest.NewRequest(http.MethodPost, "/apikey", nil)
	goodReq = mux.SetURLVars(goodReq, map[string]string{})

	goodReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	noBodyReq := httptest.NewRequest(http.MethodPost, "/apikey", nil)
	noBodyReq = mux.SetURLVars(noBodyReq, map[string]string{})

	testCases := []struct {
		Name          string
		Request       *http.Request
		Payload       *APIKeyPayload
		ExpectedError bool
	}{
		{
			Name:    "Test Case: 1 - correct body, no error",
			Request: goodReq, Payload: &APIKeyPayload{
				Email: "rafael@gmail.com",
			}, ExpectedError: false,
		},
		{Name: "Test Case: 2 - no body, error found", Request: noBodyReq, Payload: nil, ExpectedError: true},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		decodedPayload, err := decodeAPIKeyRequest(tc.Request)
		if tc.ExpectedError {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, *tc.Payload, decodedPayload)
		}
	}
}
