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
