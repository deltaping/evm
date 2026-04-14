package bankerc20

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/evm/precompiles/bankerc20"
	"github.com/cosmos/evm/testutil/integration/evm/factory"
	"github.com/cosmos/evm/testutil/integration/evm/grpc"
	"github.com/cosmos/evm/testutil/integration/evm/network"
	"github.com/cosmos/evm/testutil/integration/evm/utils"
	testkeyring "github.com/cosmos/evm/testutil/keyring"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

type PrecompileTestSuite struct {
	suite.Suite

	create                  network.CreateEvmApp
	bondDenom, tokenDenom   string
	cosmosEVMAddr, xmplAddr common.Address

	network     *network.UnitTestNetwork
	factory     factory.TxFactory
	grpcHandler grpc.Handler
	keyring     testkeyring.Keyring

	precompile *bankerc20.Precompile
}

func NewPrecompileTestSuite(create network.CreateEvmApp) *PrecompileTestSuite {
	return &PrecompileTestSuite{
		create: create,
	}
}

func (s *PrecompileTestSuite) SetupTest() sdk.Context {
	s.tokenDenom = "xmpl"

	keyring := testkeyring.New(2)
	genesis := utils.CreateGenesisWithTokenPairs(keyring)
	unitNetwork := network.NewUnitTestNetwork(
		s.create,
		network.WithPreFundedAccounts(keyring.GetAllAccAddrs()...),
		network.WithCustomGenesis(genesis),
		network.WithOtherDenoms([]string{s.tokenDenom}),
	)
	grpcHandler := grpc.NewIntegrationHandler(unitNetwork)
	txFactory := factory.New(unitNetwork, grpcHandler)

	ctx := unitNetwork.GetContext()
	sk := unitNetwork.App.GetStakingKeeperSDK()
	bondDenom, err := sk.BondDenom(ctx)
	s.Require().NoError(err, "failed to get bond denom")
	s.Require().NotEmpty(bondDenom, "bond denom cannot be empty")

	s.bondDenom = bondDenom
	s.factory = txFactory
	s.grpcHandler = grpcHandler
	s.keyring = keyring
	s.network = unitNetwork

	tokenPairID := s.network.App.GetErc20Keeper().GetTokenPairID(s.network.GetContext(), s.bondDenom)
	tokenPair, found := s.network.App.GetErc20Keeper().GetTokenPair(s.network.GetContext(), tokenPairID)
	s.Require().True(found)
	s.cosmosEVMAddr = tokenPair.GetERC20Contract()

	err = s.network.App.GetBankKeeper().MintCoins(s.network.GetContext(), minttypes.ModuleName, sdk.Coins{{Denom: s.tokenDenom, Amount: math.NewInt(1e18)}})
	s.Require().NoError(err)

	tokenPairID = s.network.App.GetErc20Keeper().GetTokenPairID(s.network.GetContext(), s.tokenDenom)
	tokenPair, found = s.network.App.GetErc20Keeper().GetTokenPair(s.network.GetContext(), tokenPairID)
	s.Require().True(found)
	s.xmplAddr = tokenPair.GetERC20Contract()

	s.precompile = s.setupPrecompile()
	return ctx
}
