package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyAuthorizedAccounts              = []byte("AuthorizedAccounts")
	DefaultAuthorizedAccounts []string = []string{""}
)

var (
	KeyDenom            = []byte("Denom")
	DefaultDenom string = "denom"
)

var (
	KeyMaxSupply            = []byte("MaxSupply")
	DefaultMaxSupply uint64 = 0
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	authorizedAccounts []string,
	denom string,
	maxSupply uint64,
) Params {
	return Params{
		AuthorizedAccounts: authorizedAccounts,
		Denom:              denom,
		MaxSupply:          maxSupply,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultAuthorizedAccounts,
		DefaultDenom,
		DefaultMaxSupply,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAuthorizedAccounts, &p.AuthorizedAccounts, validateAuthorizedAccounts),
		paramtypes.NewParamSetPair(KeyDenom, &p.Denom, validateDenom),
		paramtypes.NewParamSetPair(KeyMaxSupply, &p.MaxSupply, validateMaxSupply),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateAuthorizedAccounts(p.AuthorizedAccounts); err != nil {
		return err
	}

	if err := validateDenom(p.Denom); err != nil {
		return err
	}

	if err := validateMaxSupply(p.MaxSupply); err != nil {
		return err
	}

	return nil
}

// validateAuthorizedAccounts validates the AuthorizedAccounts param
func validateAuthorizedAccounts(v interface{}) error {
	authorizedAccounts, ok := v.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	for _, account := range authorizedAccounts {
		_, err := sdk.AccAddressFromBech32(account)
		if err != nil {
			return fmt.Errorf("invalid account address: %s", account)
		}
	}

	return nil
}

// validateDenom validates the Denom param
func validateDenom(v interface{}) error {
	denom, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	if denom == "" {
		return fmt.Errorf("denom cannot be empty")
	}

	return nil
}

// validateMaxSupply validates the MaxSupply param
func validateMaxSupply(v interface{}) error {
	maxSupply, ok := v.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	if maxSupply <= 0 {
		return fmt.Errorf("max supply must be positive")
	}

	return nil
}
