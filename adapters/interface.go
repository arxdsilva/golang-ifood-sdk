package adapters

import "io"

type Http interface {
	DoRequest(method, path string, reader io.Reader, headers map[string]string) ([]byte, int, error)
}
