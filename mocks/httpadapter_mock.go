package mocks

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type HttpClientMock struct {
	mock.Mock
}

func (h *HttpClientMock) Do(req *http.Request) (*http.Response, error) {
	args := h.Called(req)
	if res, ok := args.Get(0).(*http.Response); ok {
		return res, nil
	}
	return nil, args.Error(1)
}
