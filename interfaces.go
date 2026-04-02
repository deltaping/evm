package evm

import (
	"encoding/json"

	erc20keeper "github.com/cosmos/evm/x/erc20/keeper"
	feemarketkeeper "github.com/cosmos/evm/x/feemarket/keeper"
	precisebankkeeper "github.com/cosmos/evm/x/precisebank/keeper"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"

	storetypes "cosmossdk.io/store/types"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/mempool"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

// EvmApp defines the interface for an EVM application.
type EvmApp interface { //nolint:revive
	servertypes.ABCI

	GetBaseApp() *baseapp.BaseApp
	GetTxConfig() client.TxConfig
	AppCodec() codec.Codec
	LastCommitID() storetypes.CommitID
	LastBlockHeight() int64

	runtime.AppI
	InterfaceRegistry() types.InterfaceRegistry
	ChainID() string
	GetEVMKeeper() *evmkeeper.Keeper
	GetErc20Keeper() *erc20keeper.Keeper
	SetErc20Keeper(erc20keeper.Keeper)
	GetGovKeeper() govkeeper.Keeper
	GetSlashingKeeper() slashingkeeper.Keeper
	GetEvidenceKeeper() *evidencekeeper.Keeper
	GetBankKeeper() bankkeeper.Keeper
	GetFeeMarketKeeper() *feemarketkeeper.Keeper
	GetAccountKeeper() authkeeper.AccountKeeper
	GetAuthzKeeper() authzkeeper.Keeper
	GetDistrKeeper() distrkeeper.Keeper
	GetStakingKeeper() *stakingkeeper.Keeper
	GetStakingKeeperSDK() *stakingkeeper.Keeper
	GetMintKeeper() mintkeeper.Keeper
	GetPreciseBankKeeper() *precisebankkeeper.Keeper
	GetFeeGrantKeeper() feegrantkeeper.Keeper
	GetConsensusParamsKeeper() consensusparamkeeper.Keeper
	DefaultGenesis() map[string]json.RawMessage
	GetKey(storeKey string) *storetypes.KVStoreKey
	GetAnteHandler() sdk.AnteHandler
	MsgServiceRouter() *baseapp.MsgServiceRouter
	GetMempool() mempool.Mempool
}
