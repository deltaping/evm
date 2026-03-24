package ibc

import (
	"strings"

	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetTransferAmount returns the amount from an ICS20 FungibleTokenPacketData as a string.
func GetTransferAmount(packet channeltypes.Packet) (string, error) {
	// unmarshal packet data to obtain the sender and recipient
	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return "", errorsmod.Wrapf(errortypes.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data")
	}

	if data.Amount == "" {
		return "", errorsmod.Wrapf(errortypes.ErrInvalidCoins, "empty amount")
	}

	if _, ok := math.NewIntFromString(data.Amount); !ok {
		return "", errorsmod.Wrapf(errortypes.ErrInvalidCoins, "invalid amount")
	}

	return data.Amount, nil
}

// GetDenom returns the denomination trace from the corresponding IBC denomination. If the
// denomination is not an IBC voucher or the trace is not found, it returns an error.
func GetDenom(
	transferKeeper TransferKeeper,
	ctx sdk.Context,
	voucherDenom string,
) (transfertypes.DenomTrace, error) {
	if !strings.HasPrefix(voucherDenom, "ibc/") {
		return transfertypes.DenomTrace{}, errorsmod.Wrapf(ErrNoIBCVoucherDenom, "denom: %s", voucherDenom)
	}

	hash, err := transfertypes.ParseHexHash(voucherDenom[4:])
	if err != nil {
		return transfertypes.DenomTrace{}, err
	}

	denom, found := transferKeeper.GetDenomTrace(ctx, hash)
	if !found {
		return transfertypes.DenomTrace{}, ErrDenomNotFound
	}

	return denom, nil
}

// DeriveDecimalsFromDenom returns the number of decimals of an IBC coin
// depending on the prefix of the base denomination
func DeriveDecimalsFromDenom(baseDenom string) (uint8, error) {
	var decimals uint8
	if len(baseDenom) == 0 {
		return decimals, errorsmod.Wrapf(ErrInvalidBaseDenom, "Base denom cannot be an empty string")
	}

	switch baseDenom[0] {
	case 'u': // micro (u) -> 6 decimals
		decimals = 6
	case 'a': // atto (a) -> 18 decimals
		decimals = 18
	default:
		return decimals, errorsmod.Wrapf(
			ErrInvalidBaseDenom,
			"Should be either micro ('u[...]') or atto ('a[...]'); got: %q",
			baseDenom,
		)
	}
	return decimals, nil
}
