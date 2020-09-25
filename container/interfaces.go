package container

import "io"

type (
	HttpAdapter interface {
		DoRequest(method, path string, reader io.Reader, headers map[string]string) ([]byte, int, error)
	}
)
