package form

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	accountsPath = "/v1/organisation/accounts"
)

type d struct {
	Data interface{} `json:"data"`
}

//AccountData model as defined by form accounts API
type AccountData struct {
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	ID             string             `json:"id,omitempty"`
	OrganisationID string             `json:"organisation_id,omitempty"`
	Type           string             `json:"type,omitempty"`
	Version        *int64             `json:"version,omitempty"`

	CreatedOn  time.Time `json:"created_on,omitempty"`
	ModifiedOn time.Time `json:"modified_on,omitempty"`
}

//AccountAttributes model as defined by form accounts API
type AccountAttributes struct {
	AccountClassification   *string  `json:"account_classification,omitempty"`
	AccountMatchingOptOut   *bool    `json:"account_matching_opt_out,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	Country                 *string  `json:"country,omitempty"`
	Iban                    string   `json:"iban,omitempty"`
	JointAccount            *bool    `json:"joint_account,omitempty"`
	Name                    []string `json:"name,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Status                  *string  `json:"status,omitempty"`
	Switched                *bool    `json:"switched,omitempty"`
}

//AccountAPIClient behaves as DI container and provides methods to interact with form accounts API
type AccountAPIClient struct {
	c HTTPClient
}

//NewAccountAPIClient behaves as a construct
func NewAccountAPIClient(c HTTPClient) *AccountAPIClient {
	return &AccountAPIClient{
		c: c,
	}
}

//FetchAccountByID using GET request to "/v1/organisation/accounts/{accountID}"
func (a *AccountAPIClient) FetchAccountByID(ctx context.Context, accountID string) (*AccountData, error) {
	if err := validateAccountID(accountID); err != nil {
		return nil, ErrValidationError{Reason: "accountID isn't uuid"}
	}
	resp, err := a.c.Get(ctx, fmt.Sprintf("%s/%s", accountsPath, accountID))
	if err != nil {
		return nil, errors.Wrap(err, "GET request failed")
	}
	defer resp.Body.Close()
	switch v := resp.StatusCode; v {
	case http.StatusOK:
		b := d{Data: &AccountData{}}
		if err = json.NewDecoder(resp.Body).Decode(&b); err != nil {
			return nil, errors.Wrap(err, "failed to decode body")
		}
		return b.Data.(*AccountData), nil
	case http.StatusNotFound:
		return nil, ErrNotFound{}
	default:
		return nil, ErrUnexpectedStatusCode{StatusCode: v}
	}
}

//CreateAccountReq defines what parameters are expected in request to create new form account
type CreateAccountReq struct {
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	ID             string             `json:"id,omitempty"`
	OrganisationID string             `json:"organisation_id,omitempty"`
	Type           string             `json:"type,omitempty"`
}

//CreateAccount using POST request to const:accountsPath
func (a *AccountAPIClient) CreateAccount(ctx context.Context, req CreateAccountReq) error {
	if err := validateCreateAccountReq(req); err != nil {
		return err
	}
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(d{Data: req}); err != nil {
		return errors.Wrap(err, "failed to marshal payload to json")
	}
	resp, err := a.c.Post(ctx, accountsPath, b)
	if err != nil {
		return errors.Wrap(err, "POST request failed")
	}
	defer resp.Body.Close()

	switch v := resp.StatusCode; v {
	case http.StatusCreated:
		return nil
	case http.StatusBadRequest:
		p, err := io.ReadAll(resp.Body)
		if err != nil {
			return ErrBadRequest{Reason: "unknown"}
		}
		return ErrBadRequest{Reason: string(p)}
	case http.StatusConflict:
		return ErrConflict{Reason: "account already exists"}
	default:
		return ErrUnexpectedStatusCode{StatusCode: v}
	}
}

//DeleteAccountByID using DELETE request to "/v1/organisation/accounts/{accountID}"
func (a *AccountAPIClient) DeleteAccountByID(ctx context.Context, accountID string, version int64) error {
	if err := validateAccountID(accountID); err != nil {
		return ErrValidationError{Reason: "accountID isn't uuid"}
	}
	q := url.Values{}
	q.Set("version", fmt.Sprintf("%d", version))
	resp, err := a.c.DeleteWithQueryParams(ctx, fmt.Sprintf("%s/%s", accountsPath, accountID), q)
	if err != nil {
		return errors.Wrap(err, "DeleteWithQueryParams request failed")
	}
	defer resp.Body.Close()
	switch v := resp.StatusCode; v {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return ErrNotFound{}
	case http.StatusConflict:
		return ErrConflict{Reason: "specified version incorrect"}
	default:
		return ErrUnexpectedStatusCode{StatusCode: v}
	}
}

func validateCreateAccountReq(data CreateAccountReq) error {
	if data.Attributes == nil {
		return ErrValidationError{Reason: "Attributes property can't be empty"}
	}

	//TODO: extend to do extra validation to avoid BadRequest as much as possible
	return nil
}

func validateAccountID(accountID string) error {
	_, err := uuid.Parse(accountID)
	if err != nil {
		return ErrValidationError{Reason: "accountID isn't uuid"}
	}
	return nil
}
