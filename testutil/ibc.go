package testutil

import (
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
)

func GetVoucherDenomFromPacketData(
	data transfertypes.FungibleTokenPacketData,
	destPort string,
	destChannel string,
) string {
	// In ibc-go v8, build the voucher denom by prepending the hop path to the denom
	fullPath := destPort + "/" + destChannel + "/" + data.Denom
	trace := transfertypes.ParseDenomTrace(fullPath)
	return trace.IBCDenom()
}
