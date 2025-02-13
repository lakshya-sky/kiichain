package keeper

import (
	"fmt"
	"strings"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/kiichain/kiichain3/x/payfi/types"
)

type (
	Keeper struct {
		cdc           codec.BinaryCodec
		storeKey      sdk.StoreKey
		memKey        sdk.StoreKey
		paramstore    paramtypes.Subspace
		accountKeeper types.AccountKeeper
		bankKeeper    types.CoinTransferrer
		stakingKeeper types.StakingKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		paramstore:    ps,
		accountKeeper: accountKeeper,
		bankKeeper:    NewBankCoinTransferrer(bankKeeper),
		stakingKeeper: stakingKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) generateMerchantAddress(ctx sdk.Context, creatorAddress, businessName string) sdk.AccAddress {
	addr := buildMerchantAddress(creatorAddress, businessName)
	return addr
}

// BuildContractAddress builds an sdk account address for a contract.
func buildMerchantAddress(creatorAddress, businessName string) sdk.AccAddress {
	addressString := fmt.Sprintf("%s-%s", creatorAddress, businessName)
	addressString = strings.ToLower(addressString)
	addressString = strings.ReplaceAll(addressString, " ", "-")
	addressString = strings.ReplaceAll(addressString, "_", "-")
	addressString = strings.ReplaceAll(addressString, ".", "-")
	addressString = strings.ReplaceAll(addressString, ":", "-")
	addressString = strings.ReplaceAll(addressString, ",", "-")
	addressString = strings.ReplaceAll(addressString, "/", "-")
	addressString = strings.ReplaceAll(addressString, "\\", "-")
	addressBytes := []byte(addressString)
	return address.Module(types.ModuleName, addressBytes)[:types.MerchantAddrLen]
}

func (k Keeper) SetMerchant(ctx sdk.Context, merchant types.Merchant) {
	store := ctx.KVStore(k.storeKey)
	merchantAddr := sdk.MustAccAddressFromBech32(merchant.MerchantAddress)
	store.Set(types.MerchantKeyPrefix(merchantAddr), k.cdc.MustMarshal(&merchant))
}

func (k Keeper) GetMerchant(ctx sdk.Context, merchantAddress sdk.AccAddress) *types.Merchant {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.MerchantKeyPrefix(merchantAddress))
	if bz == nil {
		return nil
	}
	var merchant types.Merchant
	k.cdc.MustUnmarshal(bz, &merchant)
	return &merchant
}

func (k Keeper) SetPayment(ctx sdk.Context, payment types.Payment) {
	store := ctx.KVStore(k.storeKey)
	merchantAddr := sdk.MustAccAddressFromBech32(payment.MerchantAddress)
	store.Set(types.PaymentKeyPrefix(merchantAddr, payment.PaymentId), k.cdc.MustMarshal(&payment))
}

func (k Keeper) GetPayment(ctx sdk.Context, merchantAddress sdk.AccAddress, paymentId uint64) *types.Payment {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PaymentKeyPrefix(merchantAddress, paymentId))
	if bz == nil {
		return nil
	}
	var payment types.Payment
	k.cdc.MustUnmarshal(bz, &payment)
	return &payment
}

func (k Keeper) SetMerchantRevenue(ctx sdk.Context, merchantRevenue types.MerchantRevenue) {
	store := ctx.KVStore(k.storeKey)
	merchantAddr := sdk.MustAccAddressFromBech32(merchantRevenue.MerchantAddress)
	store.Set(types.MerchantRevenueKeyPrefix(merchantAddr), k.cdc.MustMarshal(&merchantRevenue))
}

func (k Keeper) GetMerchantRevenue(ctx sdk.Context, merchantAddress sdk.AccAddress) *types.MerchantRevenue {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.MerchantRevenueKeyPrefix(merchantAddress))
	if bz == nil {
		return nil
	}
	var merchantRevenue types.MerchantRevenue
	k.cdc.MustUnmarshal(bz, &merchantRevenue)
	return &merchantRevenue
}

// BankCoinTransferrer replicates the cosmos-sdk behaviour as in
// https://github.com/cosmos/cosmos-sdk/blob/v0.41.4/x/bank/keeper/msg_server.go#L26
type BankCoinTransferrer struct {
	keeper types.BankKeeper
}

func NewBankCoinTransferrer(keeper types.BankKeeper) BankCoinTransferrer {
	return BankCoinTransferrer{
		keeper: keeper,
	}
}

// TransferCoins transfers coins from source to destination account when coin send was enabled for them and the recipient
// is not in the blocked address list.
func (c BankCoinTransferrer) TransferCoins(parentCtx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amount sdk.Coins) error {
	em := sdk.NewEventManager()
	ctx := parentCtx.WithEventManager(em)
	if err := c.keeper.IsSendEnabledCoins(ctx, amount...); err != nil {
		return err
	}
	if c.keeper.BlockedAddr(toAddr) {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", toAddr.String())
	}

	sdkerr := c.keeper.SendCoins(ctx, fromAddr, toAddr, amount)
	if sdkerr != nil {
		return sdkerr
	}
	for _, e := range em.Events() {
		if e.Type == sdk.EventTypeMessage { // skip messages as we talk to the keeper directly
			continue
		}
		parentCtx.EventManager().EmitEvent(e)
	}
	return nil
}
