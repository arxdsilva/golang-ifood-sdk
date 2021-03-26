package httpadapter

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/arxdsilva/golang-ifood-sdk/mocks"
	"github.com/stretchr/testify/assert"
)

var clientMock = new(mocks.HttpClientMock)

func TestHttpAdapter_Do_Success(t *testing.T) {
	request, _ := http.NewRequest(http.MethodPost, "test/", nil)
	json := "{message: success}"
	body := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	expectedResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       body,
	}
	clientMock.On("Do", request).Once().Return(expectedResp, nil)
	adapter := New(clientMock, "test")
	resp, status, err := adapter.DoRequest(http.MethodPost, "/", nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, status)
}

func TestHttpAdapter_Do_Error(t *testing.T) {
	request, _ := http.NewRequest(http.MethodPost, "test/", nil)
	clientMock.On("Do", request).Once().Return(nil, errors.New("error"))
	adapter := New(clientMock, "test")
	resp, status, err := adapter.DoRequest(http.MethodPost, "/", nil, nil)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Equal(t, 0, status)
}

func TestNewJsonReader_NilErr(t *testing.T) {
	_, err := NewJsonReader(nil)
	assert.Equal(t, ErrorNilData, err)
}

func TestNewJsonReader_OK(t *testing.T) {
	someObject := struct {
		Name string `json:"name"`
	}{Name: "Newman"}
	reader, err := NewJsonReader(someObject)
	assert.Nil(t, err)
	assert.NotNil(t, reader)
}

func TestNewMultipartReader_NilErr(t *testing.T) {
	_, _, err := NewMultipartReader(nil)
	assert.Equal(t, ErrorNilData, err)
}

func TestNewMultipartReader_OK(t *testing.T) {
	someObject := struct {
		Name string `json:"name"`
	}{Name: "Newman"}
	reader, boudary, err := NewMultipartReader(someObject)
	assert.Nil(t, err)
	assert.NotNil(t, reader)
	assert.NotEmpty(t, boudary)
}
