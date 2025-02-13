package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/kiichain/kiichain3/x/payfi/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// RegisterMerchant handles the registration of a new merchant in the system.
// It validates the registration request, generates a unique merchant address,
// and stores the merchant information on-chain.
func (k msgServer) RegisterMerchant(goCtx context.Context, msg *types.MsgRegisterMerchant) (*types.MsgRegisterMerchantResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the msg.
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	creatorAddress := msg.CreatorAddress
	businessName := msg.BusinessName

	merchantAddress := k.generateMerchantAddress(ctx, creatorAddress, businessName)

	if k.GetMerchant(ctx, merchantAddress) != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "merchant already exists")
	}

	merchant := types.Merchant{
		MerchantAddress:     merchantAddress.String(),
		CreatorAddress:      creatorAddress,
		BusinessName:        businessName,
		BusinessAddress:     msg.BusinessAddress,
		BusinessPhone:       msg.BusinessPhone,
		BusinessEmail:       msg.BusinessEmail,
		BusinessDescription: msg.BusinessDescription,
		KycStatus:           false,
	}

	k.SetMerchant(ctx, merchant)

	return &types.MsgRegisterMerchantResponse{
		MerchantAddress: merchantAddress.String(),
	}, nil
}

// VerifyMerchant handles the KYC verification of a registered merchant.
// It updates the merchant's KYC status to true after successful verification.
func (k msgServer) VerifyMerchant(goCtx context.Context, msg *types.MsgVerifyMerchant) (*types.MsgVerifyMerchantResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the msg.
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	merchantAddress, err := sdk.AccAddressFromBech32(msg.MerchantAddress)
	if err != nil {
		return nil, err
	}

	merchant := k.GetMerchant(ctx, merchantAddress)
	if merchant == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "merchant not found")
	}

	merchant.KycStatus = true
	k.SetMerchant(ctx, *merchant)

	return &types.MsgVerifyMerchantResponse{
		MerchantAddress: merchantAddress.String(),
	}, nil
}

// PayMerchant handles payments to verified merchants.
// It validates the payment request, checks merchant verification status,
// and transfers coins from the payer to the merchant.
func (k msgServer) PayMerchant(goCtx context.Context, msg *types.MsgPayMerchant) (*types.MsgPayMerchantResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the msg.
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	merchantAddress, err := sdk.AccAddressFromBech32(msg.MerchantAddress)
	if err != nil {
		return nil, err
	}

	payerAddress, err := sdk.AccAddressFromBech32(msg.PayerAddress)
	if err != nil {
		return nil, err
	}

	merchant := k.GetMerchant(ctx, merchantAddress)
	if merchant == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "merchant not found")
	}

	if !merchant.KycStatus {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "merchant not verified")
	}

	// Send coins from payer to merchant
	err = k.bankKeeper.TransferCoins(
		ctx,
		payerAddress,
		merchantAddress,
		msg.Amount,
	)

	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to send %v from %s to %s", msg.Amount, payerAddress, merchantAddress)
	}

	return &types.MsgPayMerchantResponse{
		MerchantAddress: merchantAddress.String(),
	}, nil
}
