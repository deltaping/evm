package bankerc20

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	erc20types "github.com/cosmos/evm/x/erc20/types"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Mint mints bank coins for the denom bound to the caller contract's TokenPair
// and sends them to the specified address.
// Only callable by contracts (not EOAs). Requires a registered TokenPair.
func (p Precompile) Mint(
	ctx sdk.Context,
	contract *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	to, amount, err := ParseMintArgs(args)
	if err != nil {
		return nil, err
	}

	if amount.Sign() <= 0 {
		return nil, fmt.Errorf("mint amount must be positive")
	}

	caller := contract.Caller()
	if err := p.requireContractCaller(ctx, caller); err != nil {
		return nil, err
	}

	denom, err := p.resolveDenom(ctx, caller)
	if err != nil {
		return nil, err
	}

	coins := sdk.NewCoins(sdk.NewCoin(denom, math.NewIntFromBigInt(amount)))

	if err := p.bankKeeper.MintCoins(ctx, erc20types.ModuleName, coins); err != nil {
		return nil, fmt.Errorf("failed to mint coins: %w", err)
	}

	toAddr := sdk.AccAddress(to.Bytes())
	if err := p.bankKeeper.SendCoinsFromModuleToAccount(ctx, erc20types.ModuleName, toAddr, coins); err != nil {
		return nil, fmt.Errorf("failed to send minted coins: %w", err)
	}

	return method.Outputs.Pack(true)
}

// Burn burns bank coins for the denom bound to the caller contract's TokenPair,
// deducting them from the specified address.
// Only callable by contracts (not EOAs). Requires a registered TokenPair.
func (p Precompile) Burn(
	ctx sdk.Context,
	contract *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	from, amount, err := ParseBurnArgs(args)
	if err != nil {
		return nil, err
	}

	if amount.Sign() <= 0 {
		return nil, fmt.Errorf("burn amount must be positive")
	}

	caller := contract.Caller()
	if err := p.requireContractCaller(ctx, caller); err != nil {
		return nil, err
	}

	denom, err := p.resolveDenom(ctx, caller)
	if err != nil {
		return nil, err
	}

	coins := sdk.NewCoins(sdk.NewCoin(denom, math.NewIntFromBigInt(amount)))
	fromAddr := sdk.AccAddress(from.Bytes())

	if err := p.bankKeeper.SendCoinsFromAccountToModule(ctx, fromAddr, erc20types.ModuleName, coins); err != nil {
		return nil, fmt.Errorf("failed to send coins to module for burn: %w", err)
	}

	if err := p.bankKeeper.BurnCoins(ctx, erc20types.ModuleName, coins); err != nil {
		return nil, fmt.Errorf("failed to burn coins: %w", err)
	}

	return method.Outputs.Pack(true)
}

// Transfer transfers bank coins for the denom bound to the caller contract's
// TokenPair from one address to another.
// Only callable by contracts (not EOAs). Requires a registered TokenPair.
func (p Precompile) Transfer(
	ctx sdk.Context,
	contract *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	from, to, amount, err := ParseTransferArgs(args)
	if err != nil {
		return nil, err
	}

	if amount.Sign() <= 0 {
		return nil, fmt.Errorf("transfer amount must be positive")
	}

	caller := contract.Caller()
	if err := p.requireContractCaller(ctx, caller); err != nil {
		return nil, err
	}

	denom, err := p.resolveDenom(ctx, caller)
	if err != nil {
		return nil, err
	}

	coins := sdk.NewCoins(sdk.NewCoin(denom, math.NewIntFromBigInt(amount)))
	fromAddr := sdk.AccAddress(from.Bytes())
	toAddr := sdk.AccAddress(to.Bytes())

	if p.bankKeeper.BlockedAddr(toAddr) {
		return nil, fmt.Errorf("recipient address %s is blocked", to.Hex())
	}

	if err := p.bankKeeper.SendCoins(ctx, fromAddr, toAddr, coins); err != nil {
		return nil, fmt.Errorf("failed to transfer coins: %w", err)
	}

	return method.Outputs.Pack(true)
}

// requireContractCaller checks that the caller is a contract, not an EOA.
func (p Precompile) requireContractCaller(ctx sdk.Context, caller common.Address) error {
	if !p.evmKeeper.IsContract(ctx, caller) {
		return fmt.Errorf("caller %s is not a contract; only contracts can call this method", caller.Hex())
	}
	return nil
}
