package accounts

const (
	errorNameLength                    string = "Required length of AccountAttributes.Name is between 1 and 4"
	errorNameEmpty                     string = "No item in AccountAttributes.Name can be empty"
	errorCountryMissing                string = "AccountAttributes.Country is required"
	errorBicFormat                     string = "Required format of AccountAttributes.Bic is: " + formatBic
	errorIbanFormat                    string = "Required format of AccountAttributes.Iban is: " + formatIban
	errorAccountNumberFormat           string = "Required format of AccountAttributes.AccountNumber is: " + formatAccountNumber
	errorBankIDFormat                  string = "Required format of AccountAttributes.BankID is: " + formatBankId
	errorBankIDCodeFormat              string = "Required format of AccountAttributes.BankIDCode is: " + formatBankIdCode
	errorSecondaryIdentificationLength string = "Max length of AccountAttributes.SecondaryIdentification is 3"
	errorSecondaryIdentificationEmpty  string = "No item in AccountAttributes.SecondaryIdentification can be empty"
)
