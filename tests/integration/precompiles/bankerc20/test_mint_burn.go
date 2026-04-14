package bankerc20

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"

	"github.com/cosmos/evm/precompiles/bankerc20"
	cosmosevmutiltx "github.com/cosmos/evm/testutil/tx"
	erc20types "github.com/cosmos/evm/x/erc20/types"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (s *PrecompileTestSuite) TestMint() {
	s.SetupTest()
	method := s.precompile.Methods[bankerc20.MintMethod]

	testcases := []struct {
		name        string
		malleate    func() []interface{}
		expPass     bool
		errContains string
	}{
		{
			"fail - invalid number of arguments",
			func() []interface{} {
				return []interface{}{""}
			},
			false,
			"invalid number of arguments",
		},
		{
			"fail - zero amount",
			func() []interface{} {
				return []interface{}{
					cosmosevmutiltx.GenerateAddress(),
					big.NewInt(0),
				}
			},
			false,
			"mint amount must be positive",
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			args := tc.malleate()
			contract := vm.NewContract(
				s.keyring.GetAddr(0),
				s.precompile.Address(),
				uint256.NewInt(0),
				200_000,
				nil,
			)
			bz, err := s.precompile.Mint(
				s.network.GetContext(),
				contract,
				&method,
				args,
			)

			if tc.expPass {
				s.Require().NoError(err, "unexpected error")
				s.Require().NotNil(bz)
			} else {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errContains)
			}
		})
	}
}

func (s *PrecompileTestSuite) TestBurn() {
	s.SetupTest()
	method := s.precompile.Methods[bankerc20.BurnMethod]

	testcases := []struct {
		name        string
		malleate    func() []interface{}
		expPass     bool
		errContains string
	}{
		{
			"fail - invalid number of arguments",
			func() []interface{} {
				return []interface{}{""}
			},
			false,
			"invalid number of arguments",
		},
		{
			"fail - zero amount",
			func() []interface{} {
				return []interface{}{
					cosmosevmutiltx.GenerateAddress(),
					big.NewInt(0),
				}
			},
			false,
			"burn amount must be positive",
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			args := tc.malleate()
			contract := vm.NewContract(
				s.keyring.GetAddr(0),
				s.precompile.Address(),
				uint256.NewInt(0),
				200_000,
				nil,
			)
			bz, err := s.precompile.Burn(
				s.network.GetContext(),
				contract,
				&method,
				args,
			)

			if tc.expPass {
				s.Require().NoError(err, "unexpected error")
				s.Require().NotNil(bz)
			} else {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errContains)
			}
		})
	}
}

func (s *PrecompileTestSuite) TestBalanceOf() {
	ctx := s.SetupTest()
	method := s.precompile.Methods[bankerc20.BalanceOfMethod]

	addr := s.keyring.GetAccAddr(0)
	evmAddr := s.keyring.GetAddr(0)

	mintAmount := math.NewInt(1000000)
	coins := sdk.NewCoins(sdk.NewCoin(s.bondDenom, mintAmount))
	err := s.network.App.GetBankKeeper().MintCoins(ctx, minttypes.ModuleName, coins)
	s.Require().NoError(err)
	err = s.network.App.GetBankKeeper().SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, coins)
	s.Require().NoError(err)

	testcases := []struct {
		name        string
		malleate    func() []interface{}
		expPass     bool
		errContains string
		expBalance  *big.Int
	}{
		{
			"fail - invalid number of arguments",
			func() []interface{} {
				return []interface{}{""}
			},
			false,
			"invalid number of arguments",
			nil,
		},
		{
			"pass - returns balance for valid token and account",
			func() []interface{} {
				return []interface{}{
					s.cosmosEVMAddr,
					evmAddr,
				}
			},
			true,
			"",
			nil,
		},
		{
			"pass - returns zero for unregistered token",
			func() []interface{} {
				return []interface{}{
					cosmosevmutiltx.GenerateAddress(),
					evmAddr,
				}
			},
			true,
			"",
			big.NewInt(0),
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			args := tc.malleate()
			bz, err := s.precompile.BalanceOf(
				s.network.GetContext(),
				&method,
				args,
			)

			if tc.expPass {
				s.Require().NoError(err, "unexpected error")
				s.Require().NotNil(bz)

				out, err := method.Outputs.Unpack(bz)
				s.Require().NoError(err)
				balance, ok := out[0].(*big.Int)
				s.Require().True(ok, "expected output to be a big.Int")

				if tc.expBalance != nil {
					s.Require().True(tc.expBalance.Cmp(balance) == 0, "expected balance %s, got %s", tc.expBalance, balance)
				} else {
					s.Require().True(balance.Cmp(big.NewInt(0)) > 0, "balance should be > 0")
				}
			} else {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errContains)
			}
		})
	}
}

func (s *PrecompileTestSuite) TestMetadata() {
	s.SetupTest()
	method := s.precompile.Methods[bankerc20.MetadataMethod]

	testcases := []struct {
		name        string
		malleate    func() []interface{}
		expPass     bool
		errContains string
	}{
		{
			"fail - invalid number of arguments",
			func() []interface{} {
				return []interface{}{}
			},
			false,
			"invalid number of arguments",
		},
		{
			"pass - returns empty metadata for unregistered token",
			func() []interface{} {
				return []interface{}{
					cosmosevmutiltx.GenerateAddress(),
				}
			},
			true,
			"",
		},
		{
			"pass - returns metadata for registered token",
			func() []interface{} {
				return []interface{}{
					s.cosmosEVMAddr,
				}
			},
			true,
			"",
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			args := tc.malleate()
			bz, err := s.precompile.Metadata(
				s.network.GetContext(),
				&method,
				args,
			)

			if tc.expPass {
				s.Require().NoError(err, "unexpected error")
				s.Require().NotNil(bz)
			} else {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errContains)
			}
		})
	}
}

// Suppress unused import warnings.
var (
	_ = common.Address{}
	_ = erc20types.ModuleName
)
