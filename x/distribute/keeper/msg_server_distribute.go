package keeper

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/OptioServices/optio/x/distribute/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) Distribute(goCtx context.Context, msg *types.MsgDistribute) (*types.MsgDistributeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.IsAuthorized(ctx, msg.FromAddress) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "Unauthorized sender")
	}

	params := k.GetParams(ctx)
	startDate, err := parseDate(params.DistributionStartDate)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid distribution start date: %v", err)
	}

	currentSupply := math.NewUint(k.bankKeeper.GetSupply(ctx, params.Denom).Amount.Uint64())
	msgAmount := math.NewUint(msg.Amount)
	if currentSupply.Add(msgAmount).GT(math.NewUint(params.MaxSupply)) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "max supply exceeded")
	}

	totals, err := k.processDistributions(ctx, msg.Recipients, startDate, params)
	if err != nil {
		return nil, err
	}

	if err := k.validateDailyLimits(ctx, totals, params); err != nil {
		return nil, err
	}

	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	moduleBalance := math.NewUint(k.viewKeeper.GetBalance(ctx, moduleAddr, params.Denom).Amount.Uint64())
	if err := k.mintIfNeeded(ctx, moduleBalance, msgAmount, params.Denom); err != nil {
		return nil, err
	}

	if err := k.distributeCoins(ctx, msg.Recipients, params.Denom); err != nil {
		return nil, err
	}

	if err := k.batchUpdateDailyTotals(ctx, totals); err != nil {
		return nil, err
	}

	return &types.MsgDistributeResponse{}, nil
}

func (k msgServer) processDistributions(ctx sdk.Context, recipients []*types.Recipient, startDate time.Time, params types.Params) (map[string]uint64, error) {
	pubKey, err := parseRSAPublicKey(params.DistributionSignerPublicKey)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "invalid distribution signer public key")
	}
	today := ctx.BlockTime().Truncate(0)

	totals := make(map[string]uint64, len(recipients))
	for _, r := range recipients {
		for _, d := range r.Distributions {
			date, err := parseDate(d.DistributionDate)
			if err != nil || date.Before(startDate) || date.After(today) {
				return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid distribution date '%s'", d.DistributionDate)
			}

			if k.HasNonceBeenUsed(ctx, d.Nonce) {
				return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "nonce already used")
			}

			hashInput := fmt.Sprintf("%s:%d:%s:%d", d.DistributionDate, d.Amount, r.Address, d.Nonce)
			hash := sha256.Sum256([]byte(hashInput))
			sig, err := hex.DecodeString(d.Signature)
			if err != nil || rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], sig) != nil {
				return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "signature verification failed")
			}

			k.SetNonceUsed(ctx, d.Nonce)

			totals[d.DistributionDate] += d.Amount
		}
	}
	return totals, nil
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

func parseRSAPublicKey(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to decode PEM block")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "not an RSA public key")
	}
	return rsaPub, nil
}

func (k msgServer) validateDailyLimits(ctx context.Context, totals map[string]uint64, params types.Params) error {
	for date, total := range totals {
		limit := calculateDailyLimit(date, params)
		current, found := k.GetDailyDistributionTotal(ctx, date)
		if found && math.NewUint(total).Add(math.NewUint(current)).GT(math.NewUint(limit)) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "date '%s' exceeds daily limit", date)
		}
		if !found && total > limit {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "date '%s' exceeds daily limit", date)
		}
	}
	return nil
}

func (k msgServer) mintIfNeeded(ctx context.Context, balance, amount math.Uint, denom string) error {
	if balance.LT(amount) {
		toMint := amount.Sub(balance)
		coins := sdk.NewCoins(sdk.NewCoin(denom, math.NewIntFromUint64(toMint.Uint64())))
		return k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	}
	return nil
}

func (k msgServer) distributeCoins(ctx context.Context, recipients []*types.Recipient, denom string) error {
	for _, r := range recipients {
		acct, err := sdk.AccAddressFromBech32(r.Address)
		if err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address '%s'", r.Address)
		}
		amount := uint64(0)
		for _, d := range r.Distributions {
			amount += d.Amount
		}
		coins := sdk.NewCoins(sdk.NewCoin(denom, math.NewIntFromUint64(amount)))
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, acct, coins); err != nil {
			return err
		}
	}
	return nil
}

func (k msgServer) batchUpdateDailyTotals(ctx context.Context, totals map[string]uint64) error {
	for date, amount := range totals {
		current, found := k.GetDailyDistributionTotal(ctx, date)
		newTotal := amount
		if found {
			newTotal += current
		}
		k.SetDailyDistributionTotal(ctx, date, newTotal)
	}
	return nil
}

func calculateDailyLimit(date string, params types.Params) uint64 {
	startDate, err := parseDate(params.DistributionStartDate)
	if err != nil {
		return 0
	}
	targetDate, err := parseDate(date)
	if err != nil {
		return 0
	}

	months := monthsBetween(startDate, targetDate)
	if months < 0 {
		return 0
	}

	halvingPeriod := 1 + uint64(months)/params.MonthsInHalvingPeriod
	periodStartMonths := (halvingPeriod - 1) * params.MonthsInHalvingPeriod
	periodEndMonths := halvingPeriod * params.MonthsInHalvingPeriod
	periodStart := startDate.AddDate(0, int(periodStartMonths), 0)
	periodEnd := startDate.AddDate(0, int(periodEndMonths), -1)

	if targetDate.Before(periodStart) || !targetDate.Before(periodEnd.AddDate(0, 0, 1)) {
		return 0
	}

	periodSupply := params.MaxSupply / (1 << halvingPeriod)
	daysInPeriod := uint64(periodEnd.Sub(periodStart).Hours()/24) + 1

	return periodSupply / daysInPeriod
}

func monthsBetween(start, end time.Time) int {
	if end.Before(start) {
		return -1
	}

	years := end.Year() - start.Year()
	months := years*12 + int(end.Month()) - int(start.Month())

	// Adjust for day of month
	if end.Day() < start.Day() {
		months--
	}

	// Handle edge case where end is exactly on start's day but in a prior month
	if months < 0 {
		return 0
	}
	return months
}
