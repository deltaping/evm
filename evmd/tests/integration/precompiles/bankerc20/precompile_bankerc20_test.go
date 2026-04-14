package bankerc20

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/evm/evmd/tests/integration"
	"github.com/cosmos/evm/tests/integration/precompiles/bankerc20"
)

func TestBankERC20PrecompileTestSuite(t *testing.T) {
	s := bankerc20.NewPrecompileTestSuite(integration.CreateEvmd)
	suite.Run(t, s)
}
