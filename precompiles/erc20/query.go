package erc20

import (
	"fmt"
	"math"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// NameMethod defines the ABI method name for the ERC-20 Name
	// query.
	NameMethod = "name"
	// SymbolMethod defines the ABI method name for the ERC-20 Symbol
	// query.
	SymbolMethod = "symbol"
	// DecimalsMethod defines the ABI method name for the ERC-20 Decimals
	// query.
	DecimalsMethod = "decimals"
	// TotalSupplyMethod defines the ABI method name for the ERC-20 TotalSupply
	// query.
	TotalSupplyMethod = "totalSupply"
	// BalanceOfMethod defines the ABI method name for the ERC-20 BalanceOf
	// query.
	BalanceOfMethod = "balanceOf"
	// AllowanceMethod defines the ABI method name for the Allowance
	// query.
	AllowanceMethod = "allowance"
)

// Name returns the name of the token. If the token metadata is registered in the
// bank module, it returns its name. Otherwise, it returns execution reverted.
func (p Precompile) Name(
	ctx sdk.Context,
	_ *vm.Contract,
	_ vm.StateDB,
	method *abi.Method,
	_ []interface{},
) ([]byte, error) {
	metadata, found := p.BankKeeper.GetDenomMetaData(ctx, p.tokenPair.Denom)
	if !found {
		return nil, ConvertErrToERC20Error(fmt.Errorf("display denomination not found for denom: %q", p.tokenPair.Denom))
	}

	return method.Outputs.Pack(metadata.Name)
}

// Symbol returns the symbol of the token. If the token metadata is registered in the
// bank module, it returns its symbol. Otherwise, it returns execution reverted.
func (p Precompile) Symbol(
	ctx sdk.Context,
	_ *vm.Contract,
	_ vm.StateDB,
	method *abi.Method,
	_ []interface{},
) ([]byte, error) {
	metadata, found := p.BankKeeper.GetDenomMetaData(ctx, p.tokenPair.Denom)
	if !found {
		return nil, ConvertErrToERC20Error(fmt.Errorf("display denomination not found for denom: %q", p.tokenPair.Denom))
	}

	return method.Outputs.Pack(metadata.Symbol)
}

// Decimals returns the decimals places of the token. If the token metadata is registered in the
// bank module, it returns the display denomination exponent. Otherwise, it returns execution reverted.
func (p Precompile) Decimals(
	ctx sdk.Context,
	_ *vm.Contract,
	_ vm.StateDB,
	method *abi.Method,
	_ []interface{},
) ([]byte, error) {
	metadata, found := p.BankKeeper.GetDenomMetaData(ctx, p.tokenPair.Denom)
	if !found {
		return nil, ConvertErrToERC20Error(fmt.Errorf(
			"display denomination not found for denom: %q",
			p.tokenPair.Denom,
		))
	}

	var (
		decimals     uint32
		displayFound bool
	)
	for i := len(metadata.DenomUnits) - 1; i >= 0; i-- {
		if metadata.DenomUnits[i].Denom == metadata.Display {
			decimals = metadata.DenomUnits[i].Exponent
			displayFound = true
			break
		}
	}

	if !displayFound {
		return nil, ConvertErrToERC20Error(fmt.Errorf(
			"display denomination not found for denom: %q",
			p.tokenPair.Denom,
		))
	}

	if decimals > math.MaxUint8 {
		return nil, ConvertErrToERC20Error(fmt.Errorf(
			"uint8 overflow: invalid decimals: %d",
			decimals,
		))
	}

	return method.Outputs.Pack(uint8(decimals)) //#nosec G115 // we are checking for overflow above
}

// TotalSupply returns the amount of tokens in existence. It fetches the supply
// of the coin from the bank keeper and returns zero if not found.
func (p Precompile) TotalSupply(
	ctx sdk.Context,
	_ *vm.Contract,
	_ vm.StateDB,
	method *abi.Method,
	_ []interface{},
) ([]byte, error) {
	supply := p.BankKeeper.GetSupply(ctx, p.tokenPair.Denom)

	return method.Outputs.Pack(supply.Amount.BigInt())
}

// BalanceOf returns the amount of tokens owned by account. It fetches the balance
// of the coin from the bank keeper and returns zero if not found.
func (p Precompile) BalanceOf(
	ctx sdk.Context,
	_ *vm.Contract,
	_ vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	account, err := ParseBalanceOfArgs(args)
	if err != nil {
		return nil, err
	}

	balance := p.BankKeeper.SpendableCoin(ctx, account.Bytes(), p.tokenPair.Denom)

	return method.Outputs.Pack(balance.Amount.BigInt())
}

// Allowance returns the remaining allowance of a spender for a given owner.
func (p Precompile) Allowance(
	ctx sdk.Context,
	_ *vm.Contract,
	_ vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	owner, spender, err := ParseAllowanceArgs(args)
	if err != nil {
		return nil, err
	}

	allowance, err := p.erc20Keeper.GetAllowance(ctx, p.Address(), owner, spender)
	if err != nil {
		// NOTE: We are not returning the error here, because we want to align the behavior with
		// standard ERC20 smart contracts, which return zero if an allowance is not found.
		allowance = common.Big0
	}

	return method.Outputs.Pack(allowance)
}

