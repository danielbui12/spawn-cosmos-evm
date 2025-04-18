package keeper

import (
	"github.com/rollchains/demospawnevmchain/x/nsibc/types"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/log"

	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	portkeeper "github.com/cosmos/ibc-go/v8/modules/core/05-port/keeper"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	nameservicekeeper "github.com/rollchains/demospawnevmchain/x/nameservice/keeper"
)

// Keeper defines the module keeper.
type Keeper struct {
	storeService store.KVStoreService
	cdc          codec.BinaryCodec
	schema       collections.Schema

	ics4Wrapper  porttypes.ICS4Wrapper
	PortKeeper   *portkeeper.Keeper
	ScopedKeeper capabilitykeeper.ScopedKeeper

	NameServiceKeeper *nameservicekeeper.Keeper

	ExampleStore collections.Item[uint64]

	authority string

}

// NewKeeper creates a new Keeper instance.
func NewKeeper(
	appCodec codec.BinaryCodec,
	storeService store.KVStoreService,

	ics4Wrapper porttypes.ICS4Wrapper, // usually the IBC ChannelKeeper
	portKeeper *portkeeper.Keeper,
	scopedKeeper capabilitykeeper.ScopedKeeper,
	nsk *nameservicekeeper.Keeper,

	authority string,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:          appCodec,
		ics4Wrapper:  ics4Wrapper,
		storeService: storeService,

		PortKeeper:   portKeeper,
		ScopedKeeper: scopedKeeper,

		ExampleStore: collections.NewItem(sb, collections.NewPrefix(1), "example", collections.Uint64Value),

		NameServiceKeeper: nsk,

		authority:    authority,

	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.schema = schema

	return k
}

// WithICS4Wrapper sets the ICS4Wrapper. This function may be used after
// the keeper's creation to set the module which is above this module
// in the IBC application stack.
func (k *Keeper) WithICS4Wrapper(wrapper porttypes.ICS4Wrapper) {
	k.ics4Wrapper = wrapper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+ibcexported.ModuleName+"-"+types.ModuleName)
}
