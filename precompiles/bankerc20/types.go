package bankerc20

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	cmn "github.com/cosmos/evm/precompiles/common"
)

const (
	MintMethod     = "mint"
	BurnMethod     = "burn"
	TransferMethod = "transfer"
	BalanceOfMethod = "balanceOf"
	SupplyOfMethod  = "supplyOf"
	MetadataMethod  = "metadata"
)

// ParseMintArgs parses the call arguments for the Mint transaction.
func ParseMintArgs(args []interface{}) (common.Address, *big.Int, error) {
	if len(args) != 2 {
		return common.Address{}, nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 2, len(args))
	}

	to, ok := args[0].(common.Address)
	if !ok {
		return common.Address{}, nil, fmt.Errorf(cmn.ErrInvalidType, "to", common.Address{}, args[0])
	}

	amount, ok := args[1].(*big.Int)
	if !ok {
		return common.Address{}, nil, fmt.Errorf(cmn.ErrInvalidType, "amount", &big.Int{}, args[1])
	}

	return to, amount, nil
}

// ParseBurnArgs parses the call arguments for the Burn transaction.
func ParseBurnArgs(args []interface{}) (common.Address, *big.Int, error) {
	if len(args) != 2 {
		return common.Address{}, nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 2, len(args))
	}

	from, ok := args[0].(common.Address)
	if !ok {
		return common.Address{}, nil, fmt.Errorf(cmn.ErrInvalidType, "from", common.Address{}, args[0])
	}

	amount, ok := args[1].(*big.Int)
	if !ok {
		return common.Address{}, nil, fmt.Errorf(cmn.ErrInvalidType, "amount", &big.Int{}, args[1])
	}

	return from, amount, nil
}

// ParseTransferArgs parses the call arguments for the Transfer transaction.
func ParseTransferArgs(args []interface{}) (common.Address, common.Address, *big.Int, error) {
	if len(args) != 3 {
		return common.Address{}, common.Address{}, nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 3, len(args))
	}

	from, ok := args[0].(common.Address)
	if !ok {
		return common.Address{}, common.Address{}, nil, fmt.Errorf(cmn.ErrInvalidType, "from", common.Address{}, args[0])
	}

	to, ok := args[1].(common.Address)
	if !ok {
		return common.Address{}, common.Address{}, nil, fmt.Errorf(cmn.ErrInvalidType, "to", common.Address{}, args[1])
	}

	amount, ok := args[2].(*big.Int)
	if !ok {
		return common.Address{}, common.Address{}, nil, fmt.Errorf(cmn.ErrInvalidType, "amount", &big.Int{}, args[2])
	}

	return from, to, amount, nil
}

// ParseBalanceOfArgs parses the call arguments for the BalanceOf query.
func ParseBalanceOfArgs(args []interface{}) (common.Address, common.Address, error) {
	if len(args) != 2 {
		return common.Address{}, common.Address{}, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 2, len(args))
	}

	token, ok := args[0].(common.Address)
	if !ok {
		return common.Address{}, common.Address{}, fmt.Errorf(cmn.ErrInvalidType, "token", common.Address{}, args[0])
	}

	account, ok := args[1].(common.Address)
	if !ok {
		return common.Address{}, common.Address{}, fmt.Errorf(cmn.ErrInvalidType, "account", common.Address{}, args[1])
	}

	return token, account, nil
}

// ParseSupplyOfArgs parses the call arguments for the SupplyOf query.
func ParseSupplyOfArgs(args []interface{}) (common.Address, error) {
	if len(args) != 1 {
		return common.Address{}, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 1, len(args))
	}

	token, ok := args[0].(common.Address)
	if !ok {
		return common.Address{}, fmt.Errorf(cmn.ErrInvalidType, "token", common.Address{}, args[0])
	}

	return token, nil
}

// ParseMetadataArgs parses the call arguments for the Metadata query.
func ParseMetadataArgs(args []interface{}) (common.Address, error) {
	if len(args) != 1 {
		return common.Address{}, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 1, len(args))
	}

	token, ok := args[0].(common.Address)
	if !ok {
		return common.Address{}, fmt.Errorf(cmn.ErrInvalidType, "token", common.Address{}, args[0])
	}

	return token, nil
}
