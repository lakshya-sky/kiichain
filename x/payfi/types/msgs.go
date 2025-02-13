package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgRegisterMerchant = "register_merchant"
	TypeMsgVerifyMerchant   = "verify_merchant"
	TypeMsgPayMerchant      = "pay_merchant"
)

var _ sdk.Msg = &MsgRegisterMerchant{}

func (msg MsgRegisterMerchant) Route() string { return RouterKey }
func (msg MsgRegisterMerchant) Type() string  { return TypeMsgRegisterMerchant }
func (msg MsgRegisterMerchant) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.CreatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg MsgRegisterMerchant) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgRegisterMerchant) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.CreatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.BusinessName == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "business name cannot be empty")
	}

	if msg.BusinessAddress == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "business address cannot be empty")
	}

	if msg.BusinessPhone == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "business phone cannot be empty")
	}

	if msg.BusinessEmail == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "business email cannot be empty")
	}

	if msg.BusinessDescription == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "business description cannot be empty")
	}

	return nil
}

func (msg MsgVerifyMerchant) Route() string { return RouterKey }
func (msg MsgVerifyMerchant) Type() string  { return TypeMsgVerifyMerchant }
func (msg MsgVerifyMerchant) GetSigners() []sdk.AccAddress {
	kycAdmin, err := sdk.AccAddressFromBech32(msg.KycAdminAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{kycAdmin}
}

func (msg MsgVerifyMerchant) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgVerifyMerchant) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.MerchantAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid merchant address (%s)", err)
	}

	return nil
}

func (msg MsgPayMerchant) Route() string { return RouterKey }
func (msg MsgPayMerchant) Type() string  { return TypeMsgPayMerchant }
func (msg MsgPayMerchant) GetSigners() []sdk.AccAddress {
	payer, err := sdk.AccAddressFromBech32(msg.PayerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{payer}
}

func (msg MsgPayMerchant) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgPayMerchant) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.PayerAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid payer address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.MerchantAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid merchant address (%s)", err)
	}

	if msg.Amount.Empty() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid amount")
	}

	return nil
}
