package bankerc20

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BalanceOf returns the bank module balance of a specific token for a given account.
// The token address is resolved to a denom via the registered TokenPair.
func (p Precompile) BalanceOf(
	ctx sdk.Context,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	token, account, err := ParseBalanceOfArgs(args)
	if err != nil {
		return nil, fmt.Errorf("error in balanceOf: %s", err)
	}

	denom, err := p.resolveDenom(ctx, token)
	if err != nil {
		return method.Outputs.Pack(big.NewInt(0))
	}

	balance := p.bankKeeper.GetBalance(ctx, sdk.AccAddress(account.Bytes()), denom)
	return method.Outputs.Pack(balance.Amount.BigInt())
}

// SupplyOf returns the total supply of a specific token from the bank module.
// The token address is resolved to a denom via the registered TokenPair.
func (p Precompile) SupplyOf(
	ctx sdk.Context,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	token, err := ParseSupplyOfArgs(args)
	if err != nil {
		return nil, fmt.Errorf("error in supplyOf: %s", err)
	}

	denom, err := p.resolveDenom(ctx, token)
	if err != nil {
		return method.Outputs.Pack(big.NewInt(0))
	}

	supply := p.bankKeeper.GetSupply(ctx, denom)
	return method.Outputs.Pack(supply.Amount.BigInt())
}

// Metadata returns the bank module metadata (name, symbol, decimals) for a specific token.
// The token address is resolved to a denom via the registered TokenPair.
func (p Precompile) Metadata(
	ctx sdk.Context,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	token, err := ParseMetadataArgs(args)
	if err != nil {
		return nil, fmt.Errorf("error in metadata: %s", err)
	}

	denom, err := p.resolveDenom(ctx, token)
	if err != nil {
		return method.Outputs.Pack(false, "", "", uint8(0))
	}

	metadata, found := p.bankKeeper.GetDenomMetaData(ctx, denom)
	if !found {
		return method.Outputs.Pack(false, "", "", uint8(0))
	}

	var decimals uint8
	for _, unit := range metadata.DenomUnits {
		if unit.Exponent > uint32(decimals) { //#nosec G115
			decimals = uint8(unit.Exponent) //#nosec G115
		}
	}

	return method.Outputs.Pack(true, metadata.Name, metadata.Symbol, decimals)
}

// resolveDenom resolves an EVM address to a Cosmos denom via the registered TokenPair.
func (p Precompile) resolveDenom(ctx sdk.Context, addr common.Address) (string, error) {
	id := p.erc20Keeper.GetTokenPairID(ctx, addr.String())
	if len(id) == 0 {
		return "", fmt.Errorf("no TokenPair registered for address %s", addr.Hex())
	}
	pair, found := p.erc20Keeper.GetTokenPair(ctx, id)
	if !found {
		return "", fmt.Errorf("TokenPair not found for address %s", addr.Hex())
	}
	return pair.Denom, nil
}
