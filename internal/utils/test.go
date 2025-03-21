package utils

import (
	"bytes"
	"github.com/bubaew95/yandex-diploma/conf"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func BaseTestData(t *testing.T) (*chi.Mux, *conf.Config, *gomock.Controller) {
	config := &conf.Config{
		SecretKey: "test",
	}
	route := chi.NewRouter()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	return route, config, ctrl
}

func CreateRequest(t *testing.T, ts *httptest.Server, method string, url string, data string, token string) *http.Request {
	req, err := http.NewRequest(method, ts.URL+url, bytes.NewBufferString(data))
	require.NoError(t, err)

	if token != "" {
		req.AddCookie(&http.Cookie{
			Name:     "auth_token",
			Value:    token,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		})
	}

	return req
}

func SendUserRequest(t *testing.T, req *http.Request) *http.Response {
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}
