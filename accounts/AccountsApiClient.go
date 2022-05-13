package accounts

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"net/http"
	"net/url"
)

type accountBody struct {
	Data *AccountData `json:"data"`
}

type AccountsApiClient interface {
	GetAccount(ctx context.Context, accountId *uuid.UUID) (*AccountData, error)
	CreateAccount(ctx context.Context, organisationId *uuid.UUID, attributes *AccountAttributes) (*AccountData, error)
	DeleteAccount(ctx context.Context, accountId *uuid.UUID, version uint64) error
}

type accountsApiClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(apiUrl string) AccountsApiClient {
	return &accountsApiClient{
		BaseURL: fmt.Sprintf("%s%s", apiUrl, baseApiUrl),
		HTTPClient: &http.Client{
			Timeout: apiDefaultTimeout,
		},
	}
}

// Fetch AccountData for provided accountId
func (client *accountsApiClient) GetAccount(ctx context.Context, accountId *uuid.UUID) (*AccountData, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", client.BaseURL, url.PathEscape(accountId.String())), nil)
	if err != nil {
		return nil, err
	}

	res := AccountData{}
	if err := client.sendRequest(req, http.StatusOK, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Create new account and return it's AccountData.
// AccountAttributes.Name and AccountAttributes.Country are required.
func (client *accountsApiClient) CreateAccount(ctx context.Context, organisationId *uuid.UUID, attributes *AccountAttributes) (*AccountData, error) {
	id := uuid.NewV1()
	body := &accountBody{
		Data: &AccountData{
			Type:           "accounts",
			ID:             &id,
			OrganisationID: organisationId,
			Attributes:     attributes,
		},
	}

	b, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, client.BaseURL, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	res := AccountData{}
	if err := client.sendRequest(req, http.StatusCreated, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Delete nth version of AccountData for provided accountId
func (client *accountsApiClient) DeleteAccount(ctx context.Context, accountId *uuid.UUID, version uint64) error {

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/%s?version=%d", client.BaseURL, url.PathEscape(accountId.String()), version), nil)
	if err != nil {
		return err
	}

	if err := client.sendRequest(req, http.StatusNoContent, nil); err != nil {
		return err
	}

	return nil
}

func (client *accountsApiClient) sendRequest(request *http.Request, expectedStatus int, resultInterface interface{}) error {
	request.Header.Set("Content-Type", "application/vnd.api+json")
	request.Header.Set("Accept", "application/vnd.api+json")

	res, err := client.HTTPClient.Do(request)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != expectedStatus {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if resultInterface == nil {
		return nil
	}

	fullResponse := successResponse{
		Data: resultInterface,
	}
	if err = json.NewDecoder(res.Body).Decode(&fullResponse); err != nil {
		return err
	}

	return nil
}
