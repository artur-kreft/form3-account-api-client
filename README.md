# Artur Kreft
I am new in Go

### Solution summary:
From a user perspective I would expect that client should help to avoid "thinking" 
whenever it is possible. I used strongly typed arguments for id, currency, country 
and enumerated values, so user can easily avoid mistakes.

[Documentation](https://api-docs.form3.tech/api.html?http#organisation-accounts-create)
in "Required Account Attributes by Country" section says that for specific country, bank
and bic there should be auto-generated account number and iban. It seems that provided
test api does not implement those cases, so I did not add tests for them.

### To consider in the future:
1. Retry policy - client should handle retry policy in case random network or api issues
2. Circuit breaker - client should handle case when api is not accessible for longer time

