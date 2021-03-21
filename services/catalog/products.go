package catalog

import (
	"errors"
	"fmt"
)

type (
	Products []Product
	Product  struct {
		ID                  string   `json:"id"`
		Name                string   `json:"name"`
		Description         string   `json:"description"`
		ExternalCode        string   `json:"externalCode"`
		Image               string   `json:"image"`
		Shifts              []Shift  `json:"shifts"`
		Serving             string   `json:"serving"`
		DietaryRestrictions []string `json:"dietaryRestrictions"`
		Ean                 string   `json:"ean"`
	}
)

func (p *Product) verifyFields() (err error) {
	if p.Name == "" {
		return errors.New("Product needs a name")
	}
	if len(p.Name) > 65 {
		return errors.New("Name len is higher than 65 characters")
	}
	if len(p.Description) > 2000 {
		return errors.New("Description len is higher than 2000 characters")
	}
	if len(p.Shifts) == 0 {
		return errors.New("Product needs at least 1 shift")
	}
	serving := map[string]string{
		"SERVES_1":       "",
		"SERVES_2":       "",
		"SERVES_3":       "",
		"SERVES_4":       "",
		"NOT_APPLICABLE": "",
	}
	if _, ok := serving[p.Serving]; !ok {
		return errors.New("Serving not valid, verify docs: https://developer.ifood.com.br/reference#productcontroller_createproduct")
	}
	restrictions := map[string]string{
		"VEGETARIAN":      "",
		"VEGAN":           "",
		"ORGANIC":         "",
		"GLUTEN_FREE":     "",
		"SUGAR_FREE":      "",
		"LAC_FREE":        "",
		"ALCOHOLIC_DRINK": "",
		"NATURAL":         "",
	}
	if len(p.DietaryRestrictions) > 0 {
		for _, restriction := range p.DietaryRestrictions {
			if _, ok := restrictions[restriction]; !ok {
				return fmt.Errorf(
					"restriction '%s' does not exist in docs, see: https://developer.ifood.com.br/reference#productcontroller_createproduct",
					restriction)
			}
		}
	}
	return
}
