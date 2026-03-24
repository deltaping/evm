package ibc

import (
	cmtbytes "github.com/cometbft/cometbft/libs/bytes"

	ibctypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TransferKeeper interface {
	GetDenomTrace(ctx sdk.Context, denomHash cmtbytes.HexBytes) (ibctypes.DenomTrace, bool)
}
