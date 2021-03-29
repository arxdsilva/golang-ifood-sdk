package orders

import "errors"

var (
	// ErrOrderReferenceNotSpecified no order_id specified
	ErrOrderReferenceNotSpecified = errors.New("Order reference not specified")
	// ErrCancelCodeNotSpecified no cancel code provided
	ErrCancelCodeNotSpecified = errors.New("Order cancel code not specified")
)
