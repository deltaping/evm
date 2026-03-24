package testutil

import (
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
)

var (
	UosmoDenom    = transfertypes.ParseDenomTrace("transfer/channel-0/uosmo")
	UosmoIbcDenom = UosmoDenom.IBCDenom()

	UatomDenom    = transfertypes.ParseDenomTrace("transfer/channel-1/uatom")
	UatomIbcDenom = UatomDenom.IBCDenom()

	UAtomDenom    = transfertypes.ParseDenomTrace("transfer/channel-0/aatom")
	UAtomIbcDenom = UatomDenom.IBCDenom()

	UatomOsmoDenom    = transfertypes.ParseDenomTrace("transfer/channel-0/transfer/channel-1/uatom")
	UatomOsmoIbcDenom = UatomOsmoDenom.IBCDenom()

	AatomDenom    = transfertypes.ParseDenomTrace("transfer/channel-0/aatom")
	AatomIbcDenom = AatomDenom.IBCDenom()
)
