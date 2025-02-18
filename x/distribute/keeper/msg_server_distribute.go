package keeper

import (
	"context"
	"time"

	gomath "math"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/OptioServices/optio/x/distribute/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) Distribute(goCtx context.Context, msg *types.MsgDistribute) (*types.MsgDistributeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	authorized := k.IsAuthorized(ctx, msg.FromAddress)
	if !authorized {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "Unauthorized sender")
	}

	params := k.GetParams(ctx)
	// check if minting will exceed max supply
	currentSupply := k.bankKeeper.GetSupply(ctx, params.Denom).Amount.Uint64()
	if currentSupply+msg.Amount > params.MaxSupply {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "Max supply exceeded")
	}

	distributionStartDate, err := time.Parse("2006/01/02", params.DistributionStartDate)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "Invalid distribution start date")
	}

	distributionTotals := k.GetAllDailyDistributionTotal(ctx)
	recipientDistributionTotals := make(map[string]uint64)
	for _, recipientDistribution := range msg.Recipients {
		distributionDate, err := time.Parse("2006/01/02", recipientDistribution.DistributionDate)
		if err != nil {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "Invalid distribution date")
		}
		if distributionDate.Before(distributionStartDate) {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "Distribution date is before start date")
		}

		// Initialize or add to the total for this date
		if _, exists := recipientDistributionTotals[recipientDistribution.DistributionDate]; !exists {
			recipientDistributionTotals[recipientDistribution.DistributionDate] = 0
		}
		recipientDistributionTotals[recipientDistribution.DistributionDate] += recipientDistribution.Amount
	}

	// check if recipient distribution totals will exceed daily limits
	for date, total := range recipientDistributionTotals {
		if currentTotal, ok := distributionTotals[date]; ok {
			dailyLimit := calculateDailyLimit(date, params)
			if total+currentTotal > dailyLimit {
				return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "Distribution date '%s' total will exceed daily limit", date)
			}
		}
	}

	coin := sdk.NewCoin(params.Denom, math.NewIntFromUint64(msg.Amount))
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "Minting coins failed")
	}

	for _, recipient := range msg.Recipients {
		acct, err := sdk.AccAddressFromBech32(recipient.Address)
		if err != nil {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid recipient address")
		}

		coins := sdk.NewCoins(sdk.NewCoin(params.Denom, math.NewIntFromUint64(recipient.Amount)))
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, acct, coins)
		if err != nil {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "Sending coins failed")
		}
	}

	for date, total := range recipientDistributionTotals {
		k.UpdateDailyDistributionTotal(ctx, date, total)
	}

	return &types.MsgDistributeResponse{}, nil
}

func calculateDailyLimit(date string, params types.Params) uint64 {
	distributionStartDate, err := time.Parse("2006/01/02", params.DistributionStartDate)
	if err != nil {
		return 0
	}
	distributionDate, err := time.Parse("2006/01/02", date)
	if err != nil {
		return 0
	}

	startDate := distributionStartDate
	currentDate := distributionDate
	months := 0

	for startDate.Before(currentDate) {
		startDate = startDate.AddDate(0, 1, 0) // Move to the next month
		months++
	}

	halvingPeriod := int(gomath.Ceil(float64(months) / float64(params.MonthsInHalvingPeriod)))

	halvingPeriodStartDate := distributionStartDate.AddDate(0, halvingPeriod*int(params.MonthsInHalvingPeriod), 0)
	halvingPeriodEndDate := halvingPeriodStartDate.AddDate(0, int(params.MonthsInHalvingPeriod), 0)

	halvingPeriodSupply := params.MaxSupply / (1 << halvingPeriod) // 2 raised to the halving period
	daysInPeriod := halvingPeriodEndDate.Sub(halvingPeriodStartDate).Hours() / 24

	halvingPeriodSupplyPerDay := halvingPeriodSupply / uint64(daysInPeriod)

	return halvingPeriodSupplyPerDay
}
