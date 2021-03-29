package adapters

import "io"

// Http represents the API querier abstraction
type Http interface {
	DoRequest(method, path string, reader io.Reader, headers map[string]string) ([]byte, int, error)
}
