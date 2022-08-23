package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rafaelsanzio/go-flashscore/pkg/cache"
	"github.com/stretchr/testify/assert"
)

func mockGetCacheFunc(ctx context.Context, key string) (interface{}, error) {
	value := struct {
		Name  string
		Email string
	}{
		Name:  "Any",
		Email: "Any@email.com",
	}

	return value, nil
}

func mockGetCacheNilFunc(ctx context.Context, key string) (interface{}, error) {
	return nil, nil
}

func TestVerifyCache(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/ok", nil)

	anotherW := httptest.NewRecorder()
	anotherR := httptest.NewRequest(http.MethodGet, "/ok", nil)

	another1W := httptest.NewRecorder()
	another1R := httptest.NewRequest(http.MethodGet, "/ok", nil)

	another2W := httptest.NewRecorder()
	another2R := httptest.NewRequest(http.MethodGet, "/ok", nil)

	statusOK := http.StatusOK
	statusInternalServerError := http.StatusInternalServerError

	testCases := []struct {
		Name               string
		ResponseRecorder   *httptest.ResponseRecorder
		Request            *http.Request
		HandleCacheFunc    func(ctx context.Context, key string) (interface{}, error)
		MarshalFunc        func(v interface{}) ([]byte, error)
		WriteFunc          func(http.ResponseWriter, []byte) (int, error)
		ExpectedStatusCode *int
	}{
		{
			Name:               "Should return cache valid with status 200",
			ResponseRecorder:   w,
			Request:            r,
			HandleCacheFunc:    mockGetCacheFunc,
			MarshalFunc:        jsonMarshal,
			WriteFunc:          write,
			ExpectedStatusCode: &statusOK,
		}, {
			Name:               "Should pass without cache",
			ResponseRecorder:   anotherW,
			Request:            anotherR,
			HandleCacheFunc:    mockGetCacheNilFunc,
			MarshalFunc:        jsonMarshal,
			WriteFunc:          write,
			ExpectedStatusCode: nil,
		}, {
			Name:               "Should throw an error on Marshal func and returning 500",
			ResponseRecorder:   another1W,
			Request:            another1R,
			HandleCacheFunc:    mockGetCacheFunc,
			MarshalFunc:        fakeMarshal,
			WriteFunc:          write,
			ExpectedStatusCode: &statusInternalServerError,
		}, {
			Name:               "Should throw an error on Write func and returning 500",
			ResponseRecorder:   another2W,
			Request:            another2R,
			HandleCacheFunc:    mockGetCacheFunc,
			MarshalFunc:        jsonMarshal,
			WriteFunc:          fakeWrite,
			ExpectedStatusCode: &statusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Logf(tc.Name)

		cache.SetStore(cache.MockCacheStore{
			GetFunc: tc.HandleCacheFunc,
		})
		defer cache.SetStore(nil)

		jsonMarshal = tc.MarshalFunc
		defer restoreMarshal(jsonMarshal)

		write = tc.WriteFunc
		defer restoreWrite(write)

		VerifyCache(tc.ResponseRecorder, tc.Request)

		if tc.ExpectedStatusCode != nil {
			assert.Equal(t, *tc.ExpectedStatusCode, tc.ResponseRecorder.Code)
		}
	}

}
