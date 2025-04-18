package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/rollchains/demospawnevmchain/x/nameservice/types"
)

var _ types.QueryServer = Querier{}

type Querier struct {
	Keeper
}

func NewQuerier(keeper Keeper) Querier {
	return Querier{Keeper: keeper}
}

func (k Querier) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	p, err := k.Keeper.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryParamsResponse{Params: &p}, nil
}

// ResolveName implements types.QueryServer.
func (k Querier) ResolveName(goCtx context.Context, req *types.QueryResolveNameRequest) (*types.QueryResolveNameResponse, error) {
	v, err := k.Keeper.NameMapping.Get(goCtx, req.Wallet)
	if err != nil {
		return nil, err
	}

	return &types.QueryResolveNameResponse{
		Name: v,
	}, nil
}

// ResolveWallet implements types.QueryServer.
func (k Querier) ResolveWallet(goCtx context.Context, req *types.QueryResolveWalletRequest) (*types.QueryResolveWalletResponse, error) {
	// create a way to iterate over all the name mappings.
	iter, err := k.Keeper.NameMapping.Iterate(goCtx, nil)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		// get the value (name)
		v, err := iter.Value()
		if err != nil {
			return nil, err
		}

		// if current name matches the requested name,
		// return the wallet address for the name
		if v == req.Name {
			walletAddr, err := iter.Key()
			if err != nil {
				return nil, err
			}

			return &types.QueryResolveWalletResponse{
				Wallet: walletAddr,
			}, nil
		}
	}

	return nil, fmt.Errorf("wallet not found for name %s", req.Name)
}
