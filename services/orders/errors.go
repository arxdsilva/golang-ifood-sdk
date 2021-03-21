package orders

import "errors"

var (
	ErrOrderReferenceNotSpecified = errors.New("Order reference not specified")
	ErrCancelCodeNotSpecified     = errors.New("Order cancel code not specified")
)
