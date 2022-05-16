package accounts

import (
	"context"
	"flag"
	"github.com/biter777/countries"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var apiUrlFlag = flag.String("api_url", "http://localhost:8080/v1", "url of working api")

func cleanup(t *testing.T, client AccountsApiClient, accountData *AccountData) {
	if accountData != nil {
		err := client.DeleteAccount(context.Background(), accountData.ID, *accountData.Version)
		assert.NoError(t, err)
	}
}

func TestAccountsClient_ShouldFailForWrongApiUrl(t *testing.T) {
	client := NewClient("wrong_api_url")
	id := uuid.NewV1()

	_, err := client.GetAccount(context.Background(), &id)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "unsupported protocol scheme")
}

func TestAccountsClient_GetShouldFailForMissingAccount(t *testing.T) {
	c := NewClient(*apiUrlFlag)
	id := uuid.NewV1()

	_, errGet := c.GetAccount(context.Background(), &id)
	assert.Error(t, errGet)
	assert.ErrorContains(t, errGet, id.String())
}

func TestAccountsClient_GetShouldReturnCorrectAttributes(t *testing.T) {
	client := NewClient(*apiUrlFlag)

	organisationId := []uuid.UUID{uuid.NewV1(), uuid.NewV1()}
	name := [][]string{{"a"}, {"Jan", "Kowalski"}}
	bankId := []string{"400300", "5696993"}
	bankIdCode := []string{"GBDSC", "PLOIHFD"}
	bic := []string{"NWBKGB22", "KJILDJ99"}
	joint := []bool{false, true}
	matching := []bool{false, true}
	accountNumber := []string{"40030041426819", "SKJINJKS98854"}
	iban := []string{"GB11NWBK40030041426819", "PL99999999999999999999"}
	secondaryId := []string{"A1B2C3D4", "34GJ98FD"}
	country := []countries.CountryCode{countries.ByName("PL"), countries.ByName("DE")}
	ccy := []countries.CurrencyCode{countries.CurrencyAED, countries.CurrencyBIF}
	classification := []EnumClassification{Business, Personal}
	status := []EnumStatus{Pending, Confirmed}
	alternative := [][]string{{"a", "a", "a"}, {"Mike", "Jordan"}}
	switched := []bool{false, true}
	attributes := [2]AccountAttributes{}

	for i := 0; i < len(attributes); i++ {
		attributes[i] = AccountAttributes{
			Name:                    name[i],
			Country:                 &country[i],
			BaseCurrency:            &ccy[i],
			AccountClassification:   &classification[i],
			BankID:                  bankId[i],
			BankIDCode:              bankIdCode[i],
			Bic:                     bic[i],
			SecondaryIdentification: secondaryId[i],
			JointAccount:            &joint[i],
			AccountMatchingOptOut:   &matching[i],
			AccountNumber:           accountNumber[i],
			Iban:                    iban[i],
			AlternativeNames:        alternative[i],
			Status:                  &status[i],
			Switched:                &switched[i],
		}
	}

	tests := []struct {
		name           string
		organisationId *uuid.UUID
		attributes     *AccountAttributes
	}{
		{
			name:           "Returned attibutes should be equal",
			organisationId: &organisationId[0],
			attributes:     &attributes[0],
		},
		{
			name:           "Returned attibutes should be equal",
			organisationId: &organisationId[1],
			attributes:     &attributes[1],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resCreate, errCreate := client.CreateAccount(context.Background(), tt.organisationId, tt.attributes)
			defer cleanup(t, client, resCreate)

			assert.NoError(t, errCreate)

			resGet, errGet := client.GetAccount(context.Background(), resCreate.ID)
			assert.NotNil(t, resGet)
			assert.NoError(t, errGet)
			assert.EqualValues(t, resCreate.ID, resGet.ID)
			assert.EqualValues(t, tt.attributes, resGet.Attributes)
			assert.EqualValues(t, tt.organisationId, resGet.OrganisationID)
			assert.EqualValues(t, 0, *resGet.Version)
		})
	}

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
			expectedError: "name in body is required",
		},
		{
			name:          "Should fail when empty name",
			attributes:    &AccountAttributes{Name: []string{}},
			expectedError: "name in body is required",
		},
		{
			name:          "Should fail when name is too long",
			attributes:    &AccountAttributes{Name: []string{"a", "b", "c", "d", "e"}},
			expectedError: "name in body should have at most 4 items",
		},
		{
			name:          "Should fail when first name is empy",
			attributes:    &AccountAttributes{Name: []string{"", "a", "b", "c"}},
			expectedError: "should be at least 1 chars long",
		},
		{
			name:          "Should fail when second name is empy",
			attributes:    &AccountAttributes{Name: []string{"a", "", "b", "c"}},
			expectedError: "should be at least 1 chars long",
		},
		{
			name:          "Should fail when third name is empy",
			attributes:    &AccountAttributes{Name: []string{"a", "b", "", "c"}},
			expectedError: "should be at least 1 chars long",
		},
		{
			name:          "Should fail when fourth name is empy",
			attributes:    &AccountAttributes{Name: []string{"a", "b", "v", ""}},
			expectedError: "should be at least 1 chars long",
		},
		{
			name:          "Should fail when country is missing",
			attributes:    &AccountAttributes{Name: []string{"a"}},
			expectedError: "country in body is required",
		},
		{
			name:          "Should fail when BIC has incorrect format",
			attributes:    &AccountAttributes{Name: []string{"a"}, Country: &country, Bic: "nwbkgb22"},
			expectedError: "bic in body should match '^([A-Z]{6}[A-Z0-9]{2}|[A-Z]{6}[A-Z0-9]{5})$'",
		},

		{
			name:          "Should fail when BIC has incorrect format",
			attributes:    &AccountAttributes{Name: []string{"a"}, Country: &country, Bic: "0WBKGB22"},
			expectedError: "bic in body should match '^([A-Z]{6}[A-Z0-9]{2}|[A-Z]{6}[A-Z0-9]{5})$'",
		},
		{
			name:          "Should fail when BIC has incorrect format",
			attributes:    &AccountAttributes{Name: []string{"a"}, Country: &country, Bic: "jlkjiogfdhgsd"},
			expectedError: "bic in body should match '^([A-Z]{6}[A-Z0-9]{2}|[A-Z]{6}[A-Z0-9]{5})$'",
		},
		{
			name:          "Should fail when Bic has incorrect format",
			attributes:    &AccountAttributes{Name: []string{"a"}, Country: &country, Bic: "NWBKGB22A"},
			expectedError: "bic in body should match '^([A-Z]{6}[A-Z0-9]{2}|[A-Z]{6}[A-Z0-9]{5})$'",
		},
		{
			name:          "Should fail when Iban has incorrect format",
			attributes:    &AccountAttributes{Name: []string{"a"}, Country: &country, Iban: "G111NWBK40030041426812"},
			expectedError: "iban in body should match '^[A-Z]{2}[0-9]{2}[A-Z0-9]{0,64}$'",
		},
		{
			name:          "Should fail when AccountNumber has incorrect format",
			attributes:    &AccountAttributes{Name: []string{"a"}, Country: &country, AccountNumber: "54gf"},
			expectedError: "account_number in body should match '^[A-Z0-9]{0,64}$'",
		},
		{
			name:          "Should fail when BankID has incorrect format",
			attributes:    &AccountAttributes{Name: []string{"a"}, Country: &country, BankID: "AAAAAAAA1AAAAAA1A"},
			expectedError: "bank_id in body should match '^[A-Z0-9]{0,16}$'",
		},
		{
			name:          "Should fail when BankIDCode has incorrect format",
			attributes:    &AccountAttributes{Name: []string{"a"}, Country: &country, BankIDCode: "AAAAAAAA1"},
			expectedError: "bank_id_code in body should match '^[A-Z]{0,16}$'",
		},
		{
			name:          "Should fail when AlternativeNames is too long",
			attributes:    &AccountAttributes{Name: []string{"a"}, Country: &country, AlternativeNames: []string{"a", "b", "c", "d"}},
			expectedError: "alternative_names in body should have at most 3 items",
		},
		{
			name:          "Should fail when first AlternativeName is empty",
			attributes:    &AccountAttributes{Name: []string{"a"}, Country: &country, AlternativeNames: []string{"", "b", "c"}},
			expectedError: "should be at least 1 chars long",
		},
		{
			name:          "Should fail when second AlternativeName is empty",
			attributes:    &AccountAttributes{Name: []string{"a"}, Country: &country, AlternativeNames: []string{"a", "", "c"}},
			expectedError: "should be at least 1 chars long",
		},
		{
			name:          "Should fail when third AlternativeName is empty",
			attributes:    &AccountAttributes{Name: []string{"a"}, Country: &country, AlternativeNames: []string{"a", "b", ""}},
			expectedError: "should be at least 1 chars long",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			organisation_id := uuid.NewV1()
			_, err := client.CreateAccount(context.Background(), &organisation_id, tt.attributes)
			assert.ErrorContains(t, err, tt.expectedError)
		})
	}
}

func TestAccountsClient_DeleteShouldFailForMissingAccount(t *testing.T) {
	client := NewClient(*apiUrlFlag)
	id := uuid.NewV1()

	err := client.DeleteAccount(context.Background(), &id, 0)
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

	resCreate, errCreate := client.CreateAccount(context.Background(), &organisation_id, attributes)
	defer cleanup(t, client, resCreate)

	assert.NoError(t, errCreate)

	errDelete := client.DeleteAccount(context.Background(), resCreate.ID, 1)
	assert.Error(t, errDelete)
	assert.EqualError(t, errDelete, "invalid version")
}

func TestAccountsClient_GetShouldListenToContext(t *testing.T) {
	c := NewClient(*apiUrlFlag)
	id := uuid.NewV1()

	deadlineCtx, deadlineCancel := context.WithDeadline(context.Background(), time.Now().Add(time.Minute*-10))
	defer deadlineCancel()

	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name          string
		ctx           context.Context
		expectedError string
	}{
		{
			name:          "Should fail when canceled",
			ctx:           cancelCtx,
			expectedError: "canceled",
		},
		{
			name:          "Should fail when deadline exceeded",
			ctx:           deadlineCtx,
			expectedError: "deadline exceeded",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := c.GetAccount(tt.ctx, &id)
			assert.Error(t, err)
			assert.ErrorContains(t, err, tt.expectedError)
		})
	}
}

func TestAccountsClient_DeleteShouldListenToContext(t *testing.T) {
	c := NewClient(*apiUrlFlag)
	id := uuid.NewV1()

	deadlineCtx, deadlineCancel := context.WithDeadline(context.Background(), time.Now().Add(time.Minute*-10))
	defer deadlineCancel()

	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name          string
		ctx           context.Context
		expectedError string
	}{
		{
			name:          "Should fail when canceled",
			ctx:           cancelCtx,
			expectedError: "canceled",
		},
		{
			name:          "Should fail when deadline exceeded",
			ctx:           deadlineCtx,
			expectedError: "deadline exceeded",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.DeleteAccount(tt.ctx, &id, 0)
			assert.Error(t, err)
			assert.ErrorContains(t, err, tt.expectedError)
		})
	}
}

func TestAccountsClient_CreateShouldListenToContext(t *testing.T) {
	c := NewClient(*apiUrlFlag)
	id := uuid.NewV1()

	name := []string{"a"}
	country := countries.ByName("PL")
	attributes := &AccountAttributes{
		Name:    name,
		Country: &country,
	}

	deadlineCtx, deadlineCancel := context.WithDeadline(context.Background(), time.Now().Add(time.Minute*-10))
	defer deadlineCancel()

	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name          string
		ctx           context.Context
		expectedError string
	}{
		{
			name:          "Should fail when canceled",
			ctx:           cancelCtx,
			expectedError: "canceled",
		},
		{
			name:          "Should fail when deadline exceeded",
			ctx:           deadlineCtx,
			expectedError: "deadline exceeded",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := c.CreateAccount(tt.ctx, &id, attributes)
			assert.Error(t, err)
			assert.ErrorContains(t, err, tt.expectedError)
		})
	}
}
