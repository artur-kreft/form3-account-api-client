package accounts

import "regexp"

const (
	formatBic           string = "^([A-Z]{6}[A-Z0-9]{2}|[A-Z]{6}[A-Z0-9]{5})$"
	formatBankId        string = "^[A-Z0-9]{0,16}$"
	formatBankIdCode    string = "^[A-Z]{0,16}$"
	formatAccountNumber string = "^[A-Z0-9]{0,64}$"
	formatIban          string = "^[A-Z]{2}[0-9]{2}[A-Z0-9]{0,64}$"
)

var bicRegexp, _ = regexp.Compile(formatBic)
var bankIdRegexp, _ = regexp.Compile(formatBankId)
var bankIdCodeRegexp, _ = regexp.Compile(formatBankIdCode)
var accountNumberRegexp, _ = regexp.Compile(formatAccountNumber)
var ibanRegexp, _ = regexp.Compile(formatIban)
