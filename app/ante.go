package app

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	ibcante "github.com/cosmos/ibc-go/v7/modules/core/ante"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"

	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	decorators "github.com/JunhoNetwork/junho/app/decorators"

	feeshareante "github.com/JunhoNetwork/junho/x/feeshare/ante"
	feesharekeeper "github.com/JunhoNetwork/junho/x/feeshare/keeper"
)

// HandlerOptions extends the SDK's AnteHandler options by requiring the IBC
// channel keeper and a BankKeeper with an added method for fee sharing.
type HandlerOptions struct {
	ante.HandlerOptions

	GovKeeper         *govkeeper.Keeper
	IBCKeeper         *ibckeeper.Keeper
	FeeShareKeeper    feesharekeeper.Keeper
	BankKeeperFork    feeshareante.BankKeeper
	WasmConfig        *wasmTypes.WasmConfig
	TXCounterStoreKey storetypes.StoreKey
	Cdc               codec.BinaryCodec
	StakingSubspace   paramtypes.Subspace

	BypassMinFeeMsgTypes []string
	GlobalFeeSubspace    paramtypes.Subspace
}

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	if options.AccountKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "account keeper is required for AnteHandler")
	}
	if options.BankKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "bank keeper is required for AnteHandler")
	}
	if options.SignModeHandler == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}
	if options.WasmConfig == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "wasm config is required for ante builder")
	}
	if options.TXCounterStoreKey == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "tx counter key is required for ante builder")
	}

	sigGasConsumer := options.SigGasConsumer
	if sigGasConsumer == nil {
		sigGasConsumer = ante.DefaultSigVerificationGasConsumer
	}

	anteDecorators := []sdk.AnteDecorator{
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		decorators.NewMinCommissionDecorator(options.Cdc),
		wasmkeeper.NewLimitSimulationGasDecorator(options.WasmConfig.SimulationGasLimit),
		wasmkeeper.NewCountTXDecorator(options.TXCounterStoreKey),
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		decorators.MsgFilterDecorator{},
		// decorators.NewGovPreventSpamDecorator(options.Cdc, options.GovKeeper),
		// ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		// gaiafeeante.NewFeeDecorator(options.BypassMinFeeMsgTypes, options.GlobalFeeSubspace, options.StakingSubspace, maxBypassMinFeeMsgGasUsage),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
		feeshareante.NewFeeSharePayoutDecorator(options.BankKeeperFork, options.FeeShareKeeper),
		// SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewSetPubKeyDecorator(options.AccountKeeper),
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, sigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewRedundantRelayDecorator(options.IBCKeeper),
	}

	return sdk.ChainAnteDecorators(anteDecorators...), nil
}
