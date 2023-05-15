package mint_test

import (
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/testutil/configurator"
	"github.com/stretchr/testify/require"

	"github.com/JunhoNetwork/junho/x/mint/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func TestItCreatesModuleAccountOnInitBlock(t *testing.T) {
	app, _ := simtestutil.Setup(
		configurator.NewAppConfig(
			configurator.ParamsModule(),
			configurator.AuthModule(),
			configurator.StakingModule(),
			configurator.SlashingModule(),
			configurator.TxModule(),
			configurator.ConsensusModule(),
			configurator.BankModule(),
		),
		false,
	)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.InitChain(
		abcitypes.RequestInitChain{
			AppStateBytes: []byte("{}"),
			ChainId:       "test-chain-id",
		},
	)

	acc := app.AccountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(types.ModuleName))
	require.NotNil(t, acc)
}
