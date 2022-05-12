package accounts

import (
	"encoding/json"
	"errors"
	"github.com/biter777/countries"
)

type AccountAttributes struct {
	AccountClassification   *EnumClassification     `json:"account_classification,omitempty"`
	AccountMatchingOptOut   *bool                   `json:"account_matching_opt_out,omitempty"`
	AccountNumber           string                  `json:"account_number,omitempty"`
	AlternativeNames        []string                `json:"alternative_names,omitempty"`
	BankID                  string                  `json:"bank_id,omitempty"`
	BankIDCode              string                  `json:"bank_id_code,omitempty"`
	BaseCurrency            *countries.CurrencyCode `json:"base_currency,omitempty"`
	Bic                     string                  `json:"bic,omitempty"`
	Country                 *countries.CountryCode  `json:"country,omitempty"`
	Iban                    string                  `json:"iban,omitempty"`
	JointAccount            *bool                   `json:"joint_account,omitempty"`
	Name                    []string                `json:"name,omitempty"`
	SecondaryIdentification string                  `json:"secondary_identification,omitempty"`
	Status                  *EnumStatus             `json:"status,omitempty"`
	Switched                *bool                   `json:"switched,omitempty"`
}

func (attributes *AccountAttributes) MarshalJSON() ([]byte, error) {
	type Alias AccountAttributes
	var currency *string = nil
	var country *string = nil

	if attributes.BaseCurrency != nil {
		currency = new(string)
		*currency = attributes.BaseCurrency.Alpha()
	}

	if attributes.Country != nil {
		country = new(string)
		*country = attributes.Country.Alpha2()
	}

	return json.Marshal(&struct {
		BaseCurrency *string `json:"base_currency,omitempty"`
		Country      *string `json:"country,omitempty"`
		*Alias
	}{
		BaseCurrency: currency,
		Country:      country,
		Alias:        (*Alias)(attributes),
	})
}

func (attributes *AccountAttributes) UnmarshalJSON(data []byte) error {
	type Alias AccountAttributes
	aux := &struct {
		BaseCurrency *string `json:"base_currency,omitempty"`
		Country      *string `json:"country,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(attributes),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.BaseCurrency != nil {
		currency := countries.CurrencyCodeByName(*aux.BaseCurrency)
		attributes.BaseCurrency = &currency
	}

	if aux.Country != nil {
		country := countries.ByName(*aux.Country)
		attributes.Country = &country
	}

	return nil
}

func isValid(attributes *AccountAttributes) error {
	if len(attributes.Name) == 0 || len(attributes.Name) > 4 {
		return errors.New(errorNameLength)
	}

	containsEmptyName := any(attributes.Name, func(value string) bool {
		return len(value) == 0
	})

	if containsEmptyName {
		return errors.New(errorNameEmpty)
	}

	if attributes.Country == nil {
		return errors.New(errorCountryMissing)
	}

	if len(attributes.Bic) > 0 {
		match := bicRegexp.MatchString(attributes.Bic)
		if match == false {
			return errors.New(errorBicFormat)
		}
	}

	if len(attributes.BankID) > 0 {
		match := bankIdRegexp.MatchString(attributes.BankID)
		if match == false {
			return errors.New(errorBankIDFormat)
		}
	}

	if len(attributes.BankIDCode) > 0 {
		match := bankIdCodeRegexp.MatchString(attributes.BankIDCode)
		if match == false {
			return errors.New(errorBankIDCodeFormat)
		}
	}

	if len(attributes.AccountNumber) > 0 {
		match := accountNumberRegexp.MatchString(attributes.AccountNumber)
		if match == false {
			return errors.New(errorAccountNumberFormat)
		}
	}

	if len(attributes.Iban) > 0 {
		match := ibanRegexp.MatchString(attributes.Iban)
		if match == false {
			return errors.New(errorIbanFormat)
		}
	}

	if len(attributes.AlternativeNames) > 0 {
		if len(attributes.AlternativeNames) > 3 {
			return errors.New(errorSecondaryIdentificationLength)
		}

		containsEmpty := any(attributes.AlternativeNames, func(value string) bool {
			return len(value) == 0
		})

		if containsEmpty {
			return errors.New(errorSecondaryIdentificationEmpty)
		}
	}
	return nil
}
