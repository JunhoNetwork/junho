package types

import (
	cosmossdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidateBasic runs stateless checks on the query requests
func (q QueryFeeShareRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(q.ContractAddress); err != nil {
		return cosmossdkerrors.Wrapf(err, "invalid contract address %s", q.ContractAddress)
	}

	return nil
}

// ValidateBasic runs stateless checks on the query requests
func (q QueryDeployerFeeSharesRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(q.DeployerAddress); err != nil {
		return cosmossdkerrors.Wrapf(err, "invalid deployer address %s", q.DeployerAddress)
	}

	return nil
}

// ValidateBasic runs stateless checks on the query requests
func (q QueryWithdrawerFeeSharesRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(q.WithdrawerAddress); err != nil {
		return cosmossdkerrors.Wrapf(err, "invalid withdraw address %s", q.WithdrawerAddress)
	}

	return nil
}
