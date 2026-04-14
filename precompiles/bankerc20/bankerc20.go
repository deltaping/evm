// The bankerc20 package implements a precompile that provides mint/burn/transfer
// capabilities for bank-backed ERC-20 tokens. Only callable by contracts
// (not EOAs). The caller contract's address is used to resolve the target
// denom via a registered TokenPair.

package bankerc20

import (
	"context"
	"embed"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	cmn "github.com/cosmos/evm/precompiles/common"
	evmtypes "github.com/cosmos/evm/x/vm/types"

	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

const (
	GasMint     = 50_000
	GasBurn     = 50_000
	GasTransfer = 25_000
	GasBalanceOf = 2_870
	GasSupplyOf  = 2_477
	GasMetadata  = 3_500
)

// EVMKeeper defines the interface for the EVM keeper methods used by this precompile.
type EVMKeeper interface {
	IsContract(ctx sdk.Context, addr common.Address) bool
}

// BankERC20Keeper defines the bank keeper interface needed by this precompile.
type BankERC20Keeper interface {
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetSupply(ctx context.Context, denom string) sdk.Coin
	GetDenomMetaData(ctx context.Context, denom string) (banktypes.Metadata, bool)
	SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	BlockedAddr(addr sdk.AccAddress) bool
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

var _ vm.PrecompiledContract = &Precompile{}

var (
	//go:embed abi.json
	f   embed.FS
	ABI abi.ABI
)

func init() {
	var err error
	ABI, err = cmn.LoadABI(f, "abi.json")
	if err != nil {
		panic(err)
	}
}

// Precompile defines the BankERC20 precompile.
type Precompile struct {
	cmn.Precompile

	abi.ABI
	bankKeeper  BankERC20Keeper
	erc20Keeper cmn.ERC20Keeper
	evmKeeper   EVMKeeper
}

// NewPrecompile creates a new BankERC20 Precompile instance.
func NewPrecompile(
	bankKeeper BankERC20Keeper,
	erc20Keeper cmn.ERC20Keeper,
	evmKeeper EVMKeeper,
) *Precompile {
	return &Precompile{
		Precompile: cmn.Precompile{
			KvGasConfig:          storetypes.GasConfig{},
			TransientKVGasConfig: storetypes.GasConfig{},
			ContractAddress:      common.HexToAddress(evmtypes.BankERC20PrecompileAddress),
		},
		ABI:         ABI,
		bankKeeper:  bankKeeper,
		erc20Keeper: erc20Keeper,
		evmKeeper:   evmKeeper,
	}
}

// RequiredGas calculates the precompiled contract's base gas rate.
func (p Precompile) RequiredGas(input []byte) uint64 {
	if len(input) < 4 {
		return 0
	}

	methodID := input[:4]
	method, err := p.MethodById(methodID)
	if err != nil {
		return 0
	}

	switch method.Name {
	case MintMethod:
		return GasMint
	case BurnMethod:
		return GasBurn
	case TransferMethod:
		return GasTransfer
	case BalanceOfMethod:
		return GasBalanceOf
	case SupplyOfMethod:
		return GasSupplyOf
	case MetadataMethod:
		return GasMetadata
	}

	return 0
}

func (p Precompile) Run(evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	return p.RunNativeAction(evm, contract, func(ctx sdk.Context) ([]byte, error) {
		return p.Execute(ctx, contract, readonly)
	})
}

// Execute executes the precompiled contract methods defined in the ABI.
func (p Precompile) Execute(ctx sdk.Context, contract *vm.Contract, readOnly bool) ([]byte, error) {
	method, args, err := cmn.SetupABI(p.ABI, contract, readOnly, p.IsTransaction)
	if err != nil {
		return nil, err
	}

	var bz []byte
	switch method.Name {
	case BalanceOfMethod:
		bz, err = p.BalanceOf(ctx, method, args)
	case SupplyOfMethod:
		bz, err = p.SupplyOf(ctx, method, args)
	case MetadataMethod:
		bz, err = p.Metadata(ctx, method, args)
	case MintMethod:
		bz, err = p.Mint(ctx, contract, method, args)
	case BurnMethod:
		bz, err = p.Burn(ctx, contract, method, args)
	case TransferMethod:
		bz, err = p.Transfer(ctx, contract, method, args)
	default:
		return nil, fmt.Errorf(cmn.ErrUnknownMethod, method.Name)
	}

	return bz, err
}

// IsTransaction checks if the given method name corresponds to a transaction or query.
func (Precompile) IsTransaction(method *abi.Method) bool {
	switch method.Name {
	case MintMethod, BurnMethod, TransferMethod:
		return true
	default:
		return false
	}
}
