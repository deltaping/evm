package mempool

import (
	abci "github.com/cometbft/cometbft/abci/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// CheckTxWithNonceHandler wraps a base ABCI CheckTx function and handles EVM-specific nonce gap
// errors by routing transactions to the mempool for potential future execution.
// This is needed because cosmos-sdk v0.50.x does not support SetCheckTxHandler.
func CheckTxWithNonceHandler(
	baseCheckTx func(*abci.RequestCheckTx) (*abci.ResponseCheckTx, error),
	mempool *ExperimentalEVMMempool,
) func(*abci.RequestCheckTx) (*abci.ResponseCheckTx, error) {
	return func(req *abci.RequestCheckTx) (*abci.ResponseCheckTx, error) {
		res, err := baseCheckTx(req)
		if err != nil {
			return res, err
		}

		// If CheckTx succeeded, the tx was already inserted into the mempool by BaseApp
		if res.Code == abci.CodeTypeOK {
			return res, nil
		}

		// CheckTx failed. Try to insert it as an EVM tx with an invalid nonce.
		// InsertInvalidNonce safely handles non-EVM txs (returns error, which we ignore).
		// For EVM txs, it queues them in the txpool for future execution when the nonce gap fills.
		// This replicates the behavior of the CheckTxHandler from cosmos-sdk v0.53.x.
		if insertErr := mempool.InsertInvalidNonce(req.Tx); insertErr != nil {
			// Return the txpool error (e.g., "nonce too low") instead of the ante handler error.
			// This preserves the error messages expected by clients and tests.
			return sdkerrors.ResponseCheckTxWithEvents(insertErr, uint64(res.GasWanted), uint64(res.GasUsed), nil, false), nil //#nosec G115
		}

		return res, nil
	}
}
