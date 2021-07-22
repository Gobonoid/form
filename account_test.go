package form_test

import (
	"context"
	"github.com/Gobonoid/form"
	"github.com/Gobonoid/form/client"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var (
	baseURL = os.Getenv("ACCOUNT_API_BASE_URL")
)

func TestAccountAPIClient_FetchAccountByID(t *testing.T) {

	ctx := context.Background()

	c, err := client.NewDefaultClient(baseURL)
	require.NoError(t, err)
	accounts := form.NewAccountAPIClient(c)

	tests := []struct {
		name          string
		expectErrType error
		accountID     string
		ctx           context.Context
	}{
		{
			name:          "provided accountID isn't UUID",
			expectErrType: form.ErrValidationError{},
			accountID:     "definitely-not-uuid",
			ctx:           ctx,
		},
		{
			name:          "no account exists",
			expectErrType: form.ErrNotFound{},
			accountID:     uuid.New().String(),
			ctx:           ctx,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := accounts.FetchAccountByID(ctx, tt.accountID)
			if tt.expectErrType != nil {
				assert.IsType(t, err, tt.expectErrType)
				assert.Nil(t, resp)
			}
		})
	}
}

func TestAccountAPIClient(t *testing.T) {
	ctx := context.Background()

	c, err := client.NewDefaultClient(baseURL)
	require.NoError(t, err)
	accounts := form.NewAccountAPIClient(c)

	gbCountryCode := "GB"

	tests := []struct {
		name          string
		expectErrType error
		account       form.CreateAccountReq
		ctx           context.Context
		setup         func(accountID string)
	}{
		{
			name:          "invalid account",
			expectErrType: form.ErrValidationError{},
			account:       form.CreateAccountReq{},
			ctx:           ctx,
		},
		{
			name: "duplicate transaction",
			account: form.CreateAccountReq{
				Attributes: &form.AccountAttributes{
					Country: &gbCountryCode,
					Name:    []string{"fake account"},
				},
				ID:             uuid.New().String(),
				OrganisationID: uuid.New().String(),
				Type:           "accounts",
			},
			ctx: ctx,
			setup: func(accountID string) {
				setupErr := accounts.CreateAccount(ctx, form.CreateAccountReq{
					Attributes: &form.AccountAttributes{
						Country: &gbCountryCode,
						Name:    []string{"fake account"},
					},
					ID:             accountID,
					OrganisationID: uuid.New().String(),
					Type:           "accounts",
				})
				require.NoError(t, setupErr)
			},
			expectErrType: form.ErrConflict{},
		},
		{
			name: "bad request",
			account: form.CreateAccountReq{
				Attributes: &form.AccountAttributes{
					Country: &gbCountryCode,
				},
				ID:             uuid.New().String(),
				OrganisationID: uuid.New().String(),
				Type:           "accounts",
			},
			expectErrType: form.ErrBadRequest{},
			ctx:           ctx,
		},
		{
			name: "success",
			account: form.CreateAccountReq{
				Attributes: &form.AccountAttributes{
					Country: &gbCountryCode,
					Name:    []string{"fake account"},
				},
				ID:             uuid.New().String(),
				OrganisationID: uuid.New().String(),
				Type:           "accounts",
			},
			ctx: ctx,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(tt.account.ID)
			}
			err := accounts.CreateAccount(ctx, tt.account)
			if tt.expectErrType != nil {
				assert.IsType(t, tt.expectErrType, err)
			} else {
				assert.NoError(t, err)
				fetchedAccount, err := accounts.FetchAccountByID(ctx, tt.account.ID)
				assert.NoError(t, err)
				assert.Equal(t, tt.account.ID, fetchedAccount.ID)
				assert.Equal(t, tt.account.OrganisationID, fetchedAccount.OrganisationID)
				assert.Equal(t, tt.account.Type, fetchedAccount.Type)
				assert.EqualValues(t, tt.account.Attributes, fetchedAccount.Attributes)
				assert.False(t, fetchedAccount.CreatedOn.IsZero())
				assert.False(t, fetchedAccount.ModifiedOn.IsZero())
			}
		})
	}
}

func TestAccountAPIClient_DeleteAccountByID(t *testing.T) {
	ctx := context.Background()

	c, err := client.NewDefaultClient(baseURL)
	require.NoError(t, err)
	accounts := form.NewAccountAPIClient(c)

	tests := []struct {
		name          string
		expectErrType error
		accountID     string
		ctx           context.Context
		setup         func(accountID string)
		version       int64
	}{
		{
			name:          "provided accountID isn't UUID",
			expectErrType: form.ErrValidationError{},
			accountID:     "definitely-not-uuid",
			ctx:           ctx,
		},
		{
			name:          "no account exists",
			expectErrType: form.ErrNotFound{},
			accountID:     uuid.New().String(),
			version:       0,
			ctx:           ctx,
		},
		{
			name:      "success",
			accountID: uuid.New().String(),
			version:   0,
			ctx:       ctx,
			setup: func(accountID string) {
				gbCountryCode := "GB"

				createErr := accounts.CreateAccount(ctx, form.CreateAccountReq{
					Attributes: &form.AccountAttributes{
						AccountClassification:   nil,
						AccountMatchingOptOut:   nil,
						AccountNumber:           "",
						AlternativeNames:        nil,
						BankID:                  "",
						BankIDCode:              "",
						BaseCurrency:            "",
						Bic:                     "",
						Country:                 &gbCountryCode,
						Iban:                    "",
						JointAccount:            nil,
						Name:                    []string{"fake account"},
						SecondaryIdentification: "",
						Status:                  nil,
						Switched:                nil,
					},
					ID:             accountID,
					OrganisationID: uuid.New().String(),
					Type:           "accounts",
				})
				require.NoError(t, createErr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(tt.accountID)
			}
			err = accounts.DeleteAccountByID(ctx, tt.accountID, tt.version)
			if tt.expectErrType != nil {
				assert.IsType(t, err, tt.expectErrType)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
