package catalog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProduct_verifyFields_noname(t *testing.T) {
	p := Product{}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "name")
}

func TestProduct_verifyFields_long_name(t *testing.T) {
	p := Product{
		Name: "produtoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoproduto",
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "len is higher")
}

func TestProduct_verifyFields_long_description(t *testing.T) {
	p := Product{
		Name:        "nome",
		Description: "produtoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoproduto",
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Description")
}

func TestProduct_verifyFields_no_shift(t *testing.T) {
	p := Product{
		Name: "nome",
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "shift")
}

func TestProduct_verifyFields_no_serving(t *testing.T) {
	p := Product{
		Name: "nome",
		Shifts: []Shift{
			{
				StartTime: "00:00",
				EndTime:   "23:59",
				Monday:    true,
			},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Serving")
}

func TestProduct_verifyFields_OK(t *testing.T) {
	p := Product{
		Name: "nome",
		Shifts: []Shift{
			{
				StartTime: "00:00",
				EndTime:   "23:59",
				Monday:    true,
			},
		},
		Serving: "SERVES_1",
	}
	err := p.verifyFields()
	assert.Nil(t, err)
}

func TestProduct_verifyFields_invalid_restriction(t *testing.T) {
	p := Product{
		Name: "nome",
		Shifts: []Shift{
			{
				StartTime: "00:00",
				EndTime:   "23:59",
				Monday:    true,
			},
		},
		Serving:             "SERVES_1",
		DietaryRestrictions: []string{"JAPONES"},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "restriction")
}

func TestProduct_verifyFields_OK_restriction(t *testing.T) {
	p := Product{
		Name: "nome",
		Shifts: []Shift{
			{
				StartTime: "00:00",
				EndTime:   "23:59",
				Monday:    true,
			},
		},
		Serving:             "SERVES_1",
		DietaryRestrictions: []string{"SUGAR_FREE"},
	}
	err := p.verifyFields()
	assert.Nil(t, err)
}

func TestPizza_verifyFields_ErrSizesNotSpecified(t *testing.T) {
	p := Pizza{}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrSizesNotSpecified, err)
}

func TestPizza_verifyFields_ErrCrustsNotSpecified(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "item"},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrCrustsNotSpecified, err)
}

func TestPizza_verifyFields_ErrEdgesNotSpecified(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "item"},
		},
		Crusts: []CategoryItem{
			{Name: "crust"},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrEdgesNotSpecified, err)
}

func TestPizza_verifyFields_ErrToppingsNotSpecified(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "item"},
		},
		Crusts: []CategoryItem{
			{Name: "crust"},
		},
		Edges: []CategoryItem{
			{Name: "edge"},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrToppingsNotSpecified, err)
}

func TestPizza_verifyFields_ErrShiftsNotSpecified(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "item"},
		},
		Crusts: []CategoryItem{
			{Name: "crust"},
		},
		Edges: []CategoryItem{
			{Name: "edge"},
		},
		Toppings: []CategoryItem{
			{Name: "topping"},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrShiftsNotSpecified, err)
}

func TestPizza_verifyFields_ErrSizeNameNotSpecified(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{ID: "id"},
		},
		Crusts: []CategoryItem{
			{Name: "crust"},
		},
		Edges: []CategoryItem{
			{Name: "edge"},
		},
		Toppings: []CategoryItem{
			{Name: "topping"},
		},
		Shifts: []Shift{
			{StartTime: "edge", Monday: true},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrSizeNameNotSpecified, err)
}

func TestPizza_verifyFields_ErrInvalidPizzaStatus(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "id"},
		},
		Crusts: []CategoryItem{
			{Name: "crust"},
		},
		Edges: []CategoryItem{
			{Name: "edge"},
		},
		Toppings: []CategoryItem{
			{Name: "topping"},
		},
		Shifts: []Shift{
			{StartTime: "edge", Monday: true},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidPizzaStatus, err)
}

func TestPizza_verifyFields_ErrNoAcceptedFractions(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "id", Status: "AVAILABLE"},
		},
		Crusts: []CategoryItem{
			{Name: "crust"},
		},
		Edges: []CategoryItem{
			{Name: "edge"},
		},
		Toppings: []CategoryItem{
			{Name: "topping"},
		},
		Shifts: []Shift{
			{StartTime: "edge", Monday: true},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrNoAcceptedFractions, err)
}

func TestPizza_verifyFields_ErrCrustNameNotSpecified(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "id", Status: "AVAILABLE", AcceptedFractions: []float64{1}},
		},
		Crusts: []CategoryItem{
			{ID: "crust"},
		},
		Edges: []CategoryItem{
			{Name: "edge"},
		},
		Toppings: []CategoryItem{
			{Name: "topping"},
		},
		Shifts: []Shift{
			{StartTime: "edge", Monday: true},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrCrustNameNotSpecified, err)
}

func TestPizza_verifyFields_ErrInvalidPizzaCrustStatus(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "id", Status: "AVAILABLE", AcceptedFractions: []float64{1}},
		},
		Crusts: []CategoryItem{
			{Name: "crust"},
		},
		Edges: []CategoryItem{
			{Name: "edge"},
		},
		Toppings: []CategoryItem{
			{Name: "topping"},
		},
		Shifts: []Shift{
			{StartTime: "edge", Monday: true},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidPizzaCrustStatus, err)
}

func TestPizza_verifyFields_ErrEdgeNameNotSpecified(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "id", Status: "AVAILABLE", AcceptedFractions: []float64{1}},
		},
		Crusts: []CategoryItem{
			{Name: "crust", Status: "AVAILABLE"},
		},
		Edges: []CategoryItem{
			{ID: "edge"},
		},
		Toppings: []CategoryItem{
			{Name: "topping"},
		},
		Shifts: []Shift{
			{StartTime: "edge", Monday: true},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrEdgeNameNotSpecified, err)
}

func TestPizza_verifyFields_ErrInvalidPizzaEdgeStatus(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "id", Status: "AVAILABLE", AcceptedFractions: []float64{1}},
		},
		Crusts: []CategoryItem{
			{Name: "crust", Status: "AVAILABLE"},
		},
		Edges: []CategoryItem{
			{Name: "edge"},
		},
		Toppings: []CategoryItem{
			{Name: "topping"},
		},
		Shifts: []Shift{
			{StartTime: "edge", Monday: true},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidPizzaEdgeStatus, err)
}

func TestPizza_verifyFields_ErrToppingNameNotSpecified(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "id", Status: "AVAILABLE", AcceptedFractions: []float64{1}},
		},
		Crusts: []CategoryItem{
			{Name: "crust", Status: "AVAILABLE"},
		},
		Edges: []CategoryItem{
			{Name: "edge", Status: "AVAILABLE"},
		},
		Toppings: []CategoryItem{
			{ID: "topping"},
		},
		Shifts: []Shift{
			{StartTime: "edge", Monday: true},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrToppingNameNotSpecified, err)
}

func TestPizza_verifyFields_ErrInvalidPizzaToppingStatus(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "id", Status: "AVAILABLE", AcceptedFractions: []float64{1}},
		},
		Crusts: []CategoryItem{
			{Name: "crust", Status: "AVAILABLE"},
		},
		Edges: []CategoryItem{
			{Name: "edge", Status: "AVAILABLE"},
		},
		Toppings: []CategoryItem{
			{Name: "topping"},
		},
		Shifts: []Shift{
			{StartTime: "edge", Monday: true},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidPizzaToppingStatus, err)
}

func TestPizza_verifyFields_ErrInvalidPizzaStartEndTime(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "id", Status: "AVAILABLE", AcceptedFractions: []float64{1}},
		},
		Crusts: []CategoryItem{
			{Name: "crust", Status: "AVAILABLE"},
		},
		Edges: []CategoryItem{
			{Name: "edge", Status: "AVAILABLE"},
		},
		Toppings: []CategoryItem{
			{Name: "topping", Status: "AVAILABLE"},
		},
		Shifts: []Shift{
			{StartTime: "", Monday: true},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidPizzaStartEndTime, err)
}

func TestPizza_verifyFields_ErrInvalidPizzaEndTime(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "id", Status: "AVAILABLE", AcceptedFractions: []float64{1}},
		},
		Crusts: []CategoryItem{
			{Name: "crust", Status: "AVAILABLE"},
		},
		Edges: []CategoryItem{
			{Name: "edge", Status: "AVAILABLE"},
		},
		Toppings: []CategoryItem{
			{Name: "topping", Status: "AVAILABLE"},
		},
		Shifts: []Shift{
			{StartTime: "start", EndTime: ""},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidPizzaStartEndTime, err)
}

func TestPizza_verifyFields_OK(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "id", Status: "AVAILABLE", AcceptedFractions: []float64{1}},
		},
		Crusts: []CategoryItem{
			{Name: "crust", Status: "AVAILABLE"},
		},
		Edges: []CategoryItem{
			{Name: "edge", Status: "AVAILABLE"},
		},
		Toppings: []CategoryItem{
			{Name: "topping", Status: "AVAILABLE"},
		},
		Shifts: []Shift{
			{StartTime: "start", EndTime: "end"},
		},
	}
	err := p.verifyFields()
	assert.Nil(t, err)
}
