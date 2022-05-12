package accounts

import (
	"flag"
	"github.com/biter777/countries"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var apiUrlFlag = flag.String("api_url", "http://localhost:8080/v1", "url of working api")

func cleanup(t *testing.T, client *AccountsApiClient, accountData *AccountData) {
	if accountData != nil {
		err := client.DeleteAccount(accountData.ID, *accountData.Version)
		assert.NoError(t, err)
	}
}

func TestAccountsClient_GetShouldFailForMissingAccount(t *testing.T) {
	c := NewClient(*apiUrlFlag)
	id := uuid.NewV1()

	_, errGet := c.GetAccount(id)
	assert.Error(t, errGet)
	assert.ErrorContains(t, errGet, id.String())
}

func TestAccountsClient_GetShouldReturnCorrectAttributes(t *testing.T) {
	client := NewClient(*apiUrlFlag)

	organisation_id := uuid.NewV1()
	name := []string{"a"}
	bank_id := "400300"
	bank_id_code := "GBDSC"
	bic := "NWBKGB22"
	joint := false
	matching := false
	secondary_id := "A1B2C3D4"
	country := countries.ByName("PL")
	ccy := countries.CurrencyAED
	classification := Business
	status := Pending
	alternative := []string{"a", "a", "a"}
	switched := true
	attributes := &AccountAttributes{
		Name:                    name,
		Country:                 &country,
		BaseCurrency:            &ccy,
		AccountClassification:   &classification,
		BankID:                  bank_id,
		BankIDCode:              bank_id_code,
		Bic:                     bic,
		SecondaryIdentification: secondary_id,
		JointAccount:            &joint,
		AccountMatchingOptOut:   &matching,
		AccountNumber:           "GB11NWBK40030041426819",
		Iban:                    "GB11NWBK40030041426819",
		AlternativeNames:        alternative,
		Status:                  &status,
		Switched:                &switched,
	}

	resCreate, errCreate := client.CreateAccount(organisation_id, attributes)
	defer cleanup(t, client, resCreate)

	assert.NoError(t, errCreate)

	resGet, errGet := client.GetAccount(resCreate.ID)
	assert.NotNil(t, resGet)
	assert.NoError(t, errGet)
	assert.EqualValues(t, resCreate.ID, resGet.ID)
	assert.EqualValues(t, attributes, resGet.Attributes)
	assert.EqualValues(t, organisation_id, resGet.OrganisationID)
	assert.EqualValues(t, 0, *resGet.Version)
}

func TestAccountsClient_CreateShouldReturnCorrectAttributes(t *testing.T) {

	client := NewClient(*apiUrlFlag)

	organisation_id := uuid.NewV1()
	name := []string{"a"}
	bank_id := "400300"
	bank_id_code := "GBDSC"
	bic := "NWBKGB22"
	joint := false
	matching := false
	secondary_id := "A1B2C3D4"
	country := countries.ByName("PL")
	ccy := countries.CurrencyAED
	classification := Business
	status := Pending
	alternative := []string{"a", "a", "a"}
	switched := true
	attributes := &AccountAttributes{
		Name:                    name,
		Country:                 &country,
		BaseCurrency:            &ccy,
		AccountClassification:   &classification,
		BankID:                  bank_id,
		BankIDCode:              bank_id_code,
		Bic:                     bic,
		SecondaryIdentification: secondary_id,
		JointAccount:            &joint,
		AccountMatchingOptOut:   &matching,
		AccountNumber:           "GB11NWBK40030041426819",
		Iban:                    "GB11NWBK40030041426819",
		AlternativeNames:        alternative,
		Status:                  &status,
		Switched:                &switched,
	}

	res, err := client.CreateAccount(organisation_id, attributes)
	defer cleanup(t, client, res)

	assert.NotNil(t, res)
	assert.Nil(t, err)
	assert.EqualValues(t, attributes, res.Attributes)
	assert.EqualValues(t, organisation_id, res.OrganisationID)
	assert.EqualValues(t, 0, *res.Version)
}

func TestAccountsApiClient_CreateAccountWithWrongInput(t *testing.T) {
	client := NewClient(*apiUrlFlag)
	country := countries.Poland

	tests := []struct {
		name          string
		attributes    *AccountAttributes
		expectedError string
	}{
		{
			name:          "Should fail when nil name",
			attributes:    &AccountAttributes{},
			expectedError: errorNameLength,
		},
		{
			name:          "Should fail when empty name",
			attributes:    &AccountAttributes{Name: []string{}},
			expectedError: errorNameLength,
		},
		{
			name:          "Should fail when name is too long",
			attributes:    &AccountAttributes{Name: []string{"a", "b", "c", "d", "e"}},
			expectedError: errorNameLength,
		},
		{
			name:          "Should fail when first name is empy",
			attributes:    &AccountAttributes{Name: []string{"", "a", "b", "c"}},
			expectedError: errorNameEmpty,
		},
		{
			name:          "Should fail when second name is empy",
			attributes:    &AccountAttributes{Name: []string{"a", "", "b", "c"}},
			expectedError: errorNameEmpty,
		},
		{
			name:          "Should fail when third name is empy",
			attributes:    &AccountAttributes{Name: []string{"a", "b", "", "c"}},
			expectedError: errorNameEmpty,
		},
		{
			name:          "Should fail when fourth name is empy",
			attributes:    &AccountAttributes{Name: []string{"a", "b", "v", ""}},
			expectedError: errorNameEmpty,
		},
		{
			name:          "Should fail when country is missing",
			attributes:    &AccountAttributes{Name: []string{"a"}},
			expectedError: errorCountryMissing,
		},
		{
			name:          "Should fail when BIC has incorrect format",
			attributes:    &AccountAttributes{Name: []string{"a"}, Country: &country, Bic: "nwbkgb22"},
			expectedError: errorBicFormat,
		},
		{
			name:          "Should fail when BIC has incorrect format",
			attributes:    &AccountAttributes{Name: []string{"a"}, Country: &country, Bic: "jlkjiogfdhgsd"},
			expectedError: errorBicFormat,
		},
		{
			name:          "Should fail when Bic has incorrect format",
			attributes:    &AccountAttributes{Name: []string{"a"}, Country: &country, Bic: "NWBKGB22A"},
			expectedError: errorBicFormat,
		},
		{
			name:          "Should fail when Iban has incorrect format",
			attributes:    &AccountAttributes{Name: []string{"a"}, Country: &country, Iban: "G111NWBK40030041426812"},
			expectedError: errorIbanFormat,
		},
		{
			name:          "Should fail when AccountNumber has incorrect format",
			attributes:    &AccountAttributes{Name: []string{"a"}, Country: &country, AccountNumber: "54gf"},
			expectedError: errorAccountNumberFormat,
		},
		{
			name:          "Should fail when BankID has incorrect format",
			attributes:    &AccountAttributes{Name: []string{"a"}, Country: &country, BankID: "gdgrygfjgf"},
			expectedError: errorBankIDFormat,
		},
		{
			name:          "Should fail when BankIDCode has incorrect format",
			attributes:    &AccountAttributes{Name: []string{"a"}, Country: &country, BankIDCode: "gdgrygfjgf"},
			expectedError: errorBankIDCodeFormat,
		},
		{
			name:          "Should fail when AlternativeNames is too long",
			attributes:    &AccountAttributes{Name: []string{"a"}, Country: &country, AlternativeNames: []string{"a", "b", "c", "d"}},
			expectedError: errorSecondaryIdentificationLength,
		},
		{
			name:          "Should fail when first AlternativeName is empty",
			attributes:    &AccountAttributes{Name: []string{"a"}, Country: &country, AlternativeNames: []string{"", "b", "c"}},
			expectedError: errorSecondaryIdentificationEmpty,
		},
		{
			name:          "Should fail when second AlternativeName is empty",
			attributes:    &AccountAttributes{Name: []string{"a"}, Country: &country, AlternativeNames: []string{"a", "", "c"}},
			expectedError: errorSecondaryIdentificationEmpty,
		},
		{
			name:          "Should fail when third AlternativeName is empty",
			attributes:    &AccountAttributes{Name: []string{"a"}, Country: &country, AlternativeNames: []string{"a", "b", ""}},
			expectedError: errorSecondaryIdentificationEmpty,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			organisation_id := uuid.NewV1()
			_, err := client.CreateAccount(organisation_id, tt.attributes)
			assert.EqualValues(t, tt.expectedError, err.Error())
		})
	}
}

func TestAccountsClient_DeleteShouldFailForMissingAccount(t *testing.T) {
	client := NewClient(*apiUrlFlag)
	id := uuid.NewV1()

	err := client.DeleteAccount(id, 0)
	assert.Error(t, err)
}

func TestAccountsClient_DeleteShouldFailForWrongVersion(t *testing.T) {
	client := NewClient(*apiUrlFlag)

	organisation_id := uuid.NewV1()
	name := []string{"a"}
	country := countries.ByName("PL")
	attributes := &AccountAttributes{
		Name:    name,
		Country: &country,
	}

	resCreate, errCreate := client.CreateAccount(organisation_id, attributes)
	defer cleanup(t, client, resCreate)

	assert.NoError(t, errCreate)

	errDelete := client.DeleteAccount(resCreate.ID, 1)
	assert.Error(t, errDelete)
	assert.EqualError(t, errDelete, "invalid version")
}
