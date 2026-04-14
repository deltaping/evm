package bankerc20

import (
	"github.com/cosmos/evm/precompiles/bankerc20"
)

func (s *PrecompileTestSuite) setupPrecompile() *bankerc20.Precompile {
	return bankerc20.NewPrecompile(
		s.network.App.GetBankKeeper(),
		*s.network.App.GetErc20Keeper(),
		s.network.App.GetEVMKeeper(),
	)
}
