package integration

import (
	"github.com/cosmos/cosmos-sdk/client/flags"

	dbm "github.com/cosmos/cosmos-db"

	"github.com/cosmos/evm"
	"github.com/cosmos/evm/evmd"
	srvflags "github.com/cosmos/evm/server/flags"

	clienthelpers "cosmossdk.io/client/v2/helpers"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simutils "github.com/cosmos/cosmos-sdk/testutil/sims"
)

// CreateEvmd creates an evm app for regular integration tests (non-mempool)
// This version uses a noop mempool to avoid state issues during transaction processing
func CreateEvmd(chainID string, evmChainID uint64, customBaseAppOptions ...func(*baseapp.BaseApp)) evm.EvmApp {
	defaultNodeHome, err := clienthelpers.GetNodeHomeDirectory(".evmd")
	if err != nil {
		panic(err)
	}

	db := dbm.NewMemDB()
	logger := log.NewNopLogger()
	loadLatest := true
	appOptions := NewAppOptionsWithFlagHomeAndChainID(defaultNodeHome, evmChainID)

	baseAppOptions := append(customBaseAppOptions, baseapp.SetChainID(chainID))

	return evmd.NewExampleApp(
		logger,
		db,
		nil,
		loadLatest,
		appOptions,
		baseAppOptions...,
	)
}

func NewAppOptionsWithFlagHomeAndChainID(home string, evmChainID uint64) simutils.AppOptionsMap {
	return simutils.AppOptionsMap{
		flags.FlagHome:      home,
		srvflags.EVMChainID: evmChainID,
	}
}
