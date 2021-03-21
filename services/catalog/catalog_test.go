package catalog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_verifyNewCategoryInCatalog_no_merchant(t *testing.T) {
	err := verifyNewCategoryInCatalog("", "", "", "", "")
	assert.NotNil(t, err)
	assert.Equal(t, err, ErrMerchantNotSpecified)
}

func Test_verifyNewCategoryInCatalog_no_category(t *testing.T) {
	err := verifyNewCategoryInCatalog("merchant", "", "", "", "")
	assert.NotNil(t, err)
	assert.Equal(t, err, ErrCatalogNotSpecified)
}

func Test_verifyNewCategoryInCatalog_name_too_big(t *testing.T) {
	name := "namenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamename"
	err := verifyNewCategoryInCatalog("merchant", "catalog", name, "", "")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "less than 100 characters")
}

func Test_verifyNewCategoryInCatalog_no_name(t *testing.T) {
	name := ""
	err := verifyNewCategoryInCatalog("merchant", "catalog", name, "", "")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "name on catalog")
}

func Test_verifyNewCategoryInCatalog_no_resource(t *testing.T) {
	name := "name"
	err := verifyNewCategoryInCatalog("merchant", "catalog", name, "", "")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "should be 'AVAILABLE' or 'UNAVAILABLE'")
}

func Test_verifyNewCategoryInCatalog_no_template(t *testing.T) {
	name := "name"
	err := verifyNewCategoryInCatalog("merchant", "catalog", name, "AVAILABLE", "")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "'DEFAULT' or 'PIZZA' and")
}

func Test_verifyNewCategoryInCatalog_OK(t *testing.T) {
	name := "name"
	err := verifyNewCategoryInCatalog("merchant", "catalog", name, "AVAILABLE", "DEFAULT")
	assert.Nil(t, err)
}
