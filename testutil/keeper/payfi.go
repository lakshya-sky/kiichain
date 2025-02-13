package keeper

import (
	"testing"

	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/stretchr/testify/require"

	"github.com/kiichain/kiichain3/app"
	"github.com/kiichain/kiichain3/x/payfi/keeper"
	"github.com/kiichain/kiichain3/x/payfi/types"
)

func RandomAccAddress() sdk.AccAddress {
	priv, addr := keyPubAddr()
	_ = priv
	return addr
}

func createAccountKeeper(cdc codec.Codec, kvStoreKey *sdk.KVStoreKey, subspace paramstypes.Subspace) authkeeper.AccountKeeper {
	maccPerms := map[string][]string{ // module account permissions
		authtypes.FeeCollectorName:     nil,
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		minttypes.ModuleName:           {authtypes.Minter},
		types.ModuleName:               {authtypes.Burner},
	}
	return authkeeper.NewAccountKeeper(
		cdc,
		kvStoreKey,
		subspace,
		authtypes.ProtoBaseAccount,
		maccPerms,
	)
}

func createBankKeeper(cdc codec.Codec, kvStoreKey *sdk.KVStoreKey, accountKeeper authkeeper.AccountKeeper, subspace paramstypes.Subspace, blockedAddrs map[string]bool) bankkeeper.Keeper {
	return bankkeeper.NewBaseKeeper(
		cdc,
		kvStoreKey,
		accountKeeper,
		subspace,
		blockedAddrs,
	)
}

func createPayfiKeeper(
	cdc codec.Codec,
	kvStoreKey *sdk.KVStoreKey,
	memStoreKey *sdk.MemoryStoreKey,
	subspace paramstypes.Subspace,
	authKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	stakingKeeper stakingkeeper.Keeper,
) keeper.Keeper {
	return keeper.NewKeeper(
		cdc,
		kvStoreKey,
		memStoreKey,
		subspace,
		authKeeper,
		bankKeeper,
		stakingKeeper,
	)
}

func PayfiKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	db := tmdb.NewMemDB()
	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		paramstypes.StoreKey,
		capabilitytypes.StoreKey,
		types.StoreKey,
	)
	ms := store.NewCommitMultiStore(db)
	for _, v := range keys {
		ms.MountStoreWithDB(v, sdk.StoreTypeIAVL, db)
	}
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	for _, v := range tkeys {
		ms.MountStoreWithDB(v, sdk.StoreTypeTransient, db)
	}
	memStoreKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
	for _, v := range memStoreKeys {
		ms.MountStoreWithDB(v, sdk.StoreTypeMemory, db)
	}
	require.NoError(t, ms.LoadLatestVersion())

	ctx := sdk.NewContext(ms, tmproto.Header{}, false, log.NewNopLogger())

	encodingConfig := app.MakeEncodingConfig()
	appCodec, legacyAmino := encodingConfig.Marshaler, encodingConfig.Amino

	paramsKeeper := paramskeeper.NewKeeper(
		appCodec,
		legacyAmino,
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.TStoreKey],
	)
	for _, m := range []string{
		authtypes.ModuleName,
		banktypes.ModuleName,
		stakingtypes.ModuleName,
		capabilitytypes.ModuleName,
		types.ModuleName,
	} {
		paramsKeeper.Subspace(m)
	}

	subspace := func(m string) paramstypes.Subspace {
		r, ok := paramsKeeper.GetSubspace(m)
		require.True(t, ok)
		return r
	}

	accountKeeper := createAccountKeeper(
		appCodec,
		keys[authtypes.StoreKey],
		subspace(authtypes.ModuleName),
	)

	blockedAddrs := make(map[string]bool)

	bankKeeper := createBankKeeper(
		appCodec,
		keys[banktypes.StoreKey],
		accountKeeper,
		subspace(banktypes.ModuleName),
		blockedAddrs,
	)
	bankKeeper.SetParams(ctx, banktypes.DefaultParams())

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		keys[stakingtypes.StoreKey],
		accountKeeper,
		bankKeeper,
		subspace(stakingtypes.ModuleName),
	)
	stakingKeeper.SetParams(ctx, stakingtypes.DefaultParams())

	k := createPayfiKeeper(
		appCodec,
		keys[types.StoreKey],
		memStoreKeys[types.MemStoreKey],
		subspace(types.ModuleName),
		accountKeeper,
		bankKeeper,
		stakingKeeper,
	)

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return &k, ctx
}

func PayfiKeeperWithFundedAccounts(t testing.TB, accounts []sdk.AccAddress, amount sdk.Coins) (*keeper.Keeper, sdk.Context) {

	db := tmdb.NewMemDB()
	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		paramstypes.StoreKey,
		capabilitytypes.StoreKey,
		types.StoreKey,
	)
	ms := store.NewCommitMultiStore(db)
	for _, v := range keys {
		ms.MountStoreWithDB(v, sdk.StoreTypeIAVL, db)
	}
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	for _, v := range tkeys {
		ms.MountStoreWithDB(v, sdk.StoreTypeTransient, db)
	}
	memStoreKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
	for _, v := range memStoreKeys {
		ms.MountStoreWithDB(v, sdk.StoreTypeMemory, db)
	}
	require.NoError(t, ms.LoadLatestVersion())

	ctx := sdk.NewContext(ms, tmproto.Header{}, false, log.NewNopLogger())

	encodingConfig := app.MakeEncodingConfig()
	appCodec, legacyAmino := encodingConfig.Marshaler, encodingConfig.Amino

	paramsKeeper := paramskeeper.NewKeeper(
		appCodec,
		legacyAmino,
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.TStoreKey],
	)
	for _, m := range []string{
		authtypes.ModuleName,
		banktypes.ModuleName,
		stakingtypes.ModuleName,
		capabilitytypes.ModuleName,
		types.ModuleName,
	} {
		paramsKeeper.Subspace(m)
	}

	subspace := func(m string) paramstypes.Subspace {
		r, ok := paramsKeeper.GetSubspace(m)
		require.True(t, ok)
		return r
	}

	accountKeeper := createAccountKeeper(
		appCodec,
		keys[authtypes.StoreKey],
		subspace(authtypes.ModuleName),
	)

	blockedAddrs := make(map[string]bool)

	bankKeeper := createBankKeeper(
		appCodec,
		keys[banktypes.StoreKey],
		accountKeeper,
		subspace(banktypes.ModuleName),
		blockedAddrs,
	)
	bankKeeper.SetParams(ctx, banktypes.DefaultParams())

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		keys[stakingtypes.StoreKey],
		accountKeeper,
		bankKeeper,
		subspace(stakingtypes.ModuleName),
	)
	stakingKeeper.SetParams(ctx, stakingtypes.DefaultParams())

	k := createPayfiKeeper(
		appCodec,
		keys[types.StoreKey],
		memStoreKeys[types.MemStoreKey],
		subspace(types.ModuleName),
		accountKeeper,
		bankKeeper,
		stakingKeeper,
	)

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())
	faucet := NewTestFaucet(t, ctx, bankKeeper, minttypes.ModuleName, amount...)
	for _, addr := range accounts {
		faucet.Fund(ctx, addr, amount)
	}
	return &k, ctx
}
