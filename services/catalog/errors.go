package catalog

import "errors"

var (
	// ErrMerchantNotSpecified no merchant given
	ErrMerchantNotSpecified = errors.New("merchant not specified")
	// ErrCatalogNotSpecified no catalog id given
	ErrCatalogNotSpecified = errors.New("Catalog ID was not specified")
	// ErrCategoryNotSpecified no categiry id given
	ErrCategoryNotSpecified = errors.New("Category ID was not specified")
	// ErrSizesNotSpecified no pizza size
	ErrSizesNotSpecified = errors.New("Pizza sizes were not specified")

	// ErrCrustsNotSpecified no pizza crust
	ErrCrustsNotSpecified = errors.New("Pizza crusts were not specified")
	// ErrCrustNameNotSpecified no pizza crust name
	ErrCrustNameNotSpecified = errors.New("Pizza crust name was not specified")

	// ErrEdgeNameNotSpecified no pizza edge name
	ErrEdgeNameNotSpecified = errors.New("Pizza edge name was not specified")
	// ErrEdgesNotSpecified no pizza edge
	ErrEdgesNotSpecified = errors.New("Pizza edges were not specified")

	// ErrToppingNameNotSpecified no pizza topping name
	ErrToppingNameNotSpecified = errors.New("Pizza Topping name was not specified")
	// ErrToppingsNotSpecified no pizza topping
	ErrToppingsNotSpecified = errors.New("Pizza Toppings were not specified")

	// ErrShiftsNotSpecified no shift
	ErrShiftsNotSpecified = errors.New("Pizza Shifts were not specified")

	// ErrSizeNameNotSpecified no pizza size
	ErrSizeNameNotSpecified = errors.New("Pizza size name was not specified")

	// ErrInvalidPizzaStatus no pizza status
	ErrInvalidPizzaStatus = errors.New("INVALID Pizza size status, it should be 'AVAILABLE' or 'UNAVAILABLE'")
	// ErrInvalidPizzaStartEndTime no pizza start/end time
	ErrInvalidPizzaStartEndTime = errors.New("INVALID Pizza start or end time, should be between 00:00 and 23:59")
	// ErrInvalidPizzaCrustStatus no pizza crust status
	ErrInvalidPizzaCrustStatus = errors.New("INVALID Pizza crust status, it should be 'AVAILABLE' or 'UNAVAILABLE'")
	// ErrInvalidPizzaEdgeStatus no pizza edge status
	ErrInvalidPizzaEdgeStatus = errors.New("INVALID Pizza edge status, it should be 'AVAILABLE' or 'UNAVAILABLE'")
	// ErrInvalidPizzaToppingStatus no pizza topping status
	ErrInvalidPizzaToppingStatus = errors.New("INVALID Pizza topping status, it should be 'AVAILABLE' or 'UNAVAILABLE'")

	// ErrNoAcceptedFractions no pizza fractions
	ErrNoAcceptedFractions = errors.New("Pizza needs at least one accepted fraction")

	// ErrNoProductName no product name
	ErrNoProductName = errors.New("Product needs a name")
	// ErrNoProductID no product id
	ErrNoProductID = errors.New("productID not specified")
	// ErrNoItemPrice no item price
	ErrNoItemPrice = errors.New("item needs the price value to be specified")
	// ErrNoPrice no price
	ErrNoPrice = errors.New("item needs the price struct filled")
	// ErrInvalidStatus invalid status
	ErrInvalidStatus = errors.New("INVALID status, it should be 'AVAILABLE' or 'UNAVAILABLE'")
	// ErrNoShifts no shift
	ErrNoShifts = errors.New("Item needs at least one shift")
)
