package httpadapter

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"github.com/kpango/glg"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type httpAdapter struct {
	client  HTTPClient
	baseUrl string
}

var (
	ErrorNilData = errors.New("No data to parse ")
	ErrorNilAuth = errors.New("No auth to parse ")
)

func New(client HTTPClient, baseUrl string) *httpAdapter {
	return &httpAdapter{client, baseUrl}
}

func (h *httpAdapter) DoRequest(method, path string, reader io.Reader, headers map[string]string) ([]byte, int, error) {
	request, err := http.NewRequest(method, h.baseUrl+path, reader)
	if err != nil {
		return nil, 0, err
	}
	for k, v := range headers {
		request.Header.Add(k, v)
	}
	resp, err := h.client.Do(request)
	if err != nil {
		return nil, 0, err
	}
	defer closeBodyReader(resp.Body)
	result, err := ioutil.ReadAll(resp.Body)
	return result, resp.StatusCode, err
}

func NewJsonReader(data interface{}) (io.Reader, error) {
	if data == nil {
		return nil, ErrorNilData
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		glg.Warnf("Error on makeReader marshaling json data: %e", err)
		return nil, errors.New("error on marshal data: " + err.Error())
	}
	return bytes.NewReader(jsonData), nil
}

func NewMultipartReader(data interface{}) (reader io.Reader, boundary string, err error) {
	if data == nil {
		err = ErrorNilData
		return
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		err = errors.New("error on marshal data: " + err.Error())
		glg.Warnf("Error on makeReader marshaling json data: %e", err)
		return
	}
	body := &bytes.Buffer{}
	writer, err := getWriter(body, jsonData)
	if err != nil {
		err = errors.New("error on create part data: " + err.Error())
		glg.Warnf("Error on writing metadata headers: %v", err)
		return
	}
	return bytes.NewReader(body.Bytes()), writer.Boundary(), nil
}

func getWriter(body *bytes.Buffer, data []byte) (*multipart.Writer, error) {
	writer := multipart.NewWriter(body)
	metadataHeader := textproto.MIMEHeader{}
	metadataHeader.Set("Content-Type", "application/json")
	metadataHeader.Set("Content-ID", "metadata")
	part, err := writer.CreatePart(metadataHeader)
	if err != nil {
		err = errors.New("error on create part data: " + err.Error())
		glg.Warnf("Error on writing metadata headers: %v", err)
		return nil, err
	}
	_, err = part.Write(data)
	if err != nil {
		err = errors.New("error on create part data: " + err.Error())
		glg.Warnf("Error on writing data: %v", err)
		return nil, err
	}
	if err := writer.Close(); err != nil {
		err = errors.New("error on create part data: " + err.Error())
		glg.Fatalf("Error closing multipart writer: %v", err)
		return nil, err
	}
	return writer, nil
}

func closeBodyReader(reader io.ReadCloser) {
	if err := reader.Close(); err != nil {
		glg.Warnf("Error on closeBodyReader %e", err)
	}
}
