package catalog

import "errors"

var (
	ErrMerchantNotSpecified = errors.New("merchant not specified")
	ErrCatalogNotSpecified  = errors.New("Catalog ID was not specified")
	ErrCategoryNotSpecified = errors.New("Category ID was not specified")
	ErrSizesNotSpecified    = errors.New("Pizza sizes were not specified")

	ErrCrustsNotSpecified    = errors.New("Pizza crusts were not specified")
	ErrCrustNameNotSpecified = errors.New("Pizza crust name was not specified")

	ErrEdgeNameNotSpecified = errors.New("Pizza edge name was not specified")
	ErrEdgesNotSpecified    = errors.New("Pizza edges were not specified")

	ErrToppingNameNotSpecified = errors.New("Pizza Topping name was not specified")
	ErrToppingsNotSpecified    = errors.New("Pizza Toppings were not specified")

	ErrShiftsNotSpecified = errors.New("Pizza Shifts were not specified")

	ErrSizeNameNotSpecified = errors.New("Pizza size name was not specified")

	ErrInvalidPizzaStatus        = errors.New("INVALID Pizza size status, it should be 'AVAILABLE' or 'UNAVAILABLE'")
	ErrInvalidPizzaStartEndTime  = errors.New("INVALID Pizza start or end time, should be between 00:00 and 23:59")
	ErrInvalidPizzaCrustStatus   = errors.New("INVALID Pizza crust status, it should be 'AVAILABLE' or 'UNAVAILABLE'")
	ErrInvalidPizzaEdgeStatus    = errors.New("INVALID Pizza edge status, it should be 'AVAILABLE' or 'UNAVAILABLE'")
	ErrInvalidPizzaToppingStatus = errors.New("INVALID Pizza topping status, it should be 'AVAILABLE' or 'UNAVAILABLE'")

	ErrNoAcceptedFractions = errors.New("Pizza needs at least one accepted fraction")

	ErrNoProductName = errors.New("Product needs a name")
	ErrNoProductID   = errors.New("productID not specified")
	ErrNoItemPrice   = errors.New("item needs the price value to be specified")
	ErrNoPrice       = errors.New("item needs the price struct filled")
	ErrInvalidStatus = errors.New("INVALID status, it should be 'AVAILABLE' or 'UNAVAILABLE'")
	ErrNoShifts      = errors.New("Item needs at least one shift")
)
