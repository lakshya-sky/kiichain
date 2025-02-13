package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/kiichain/kiichain3/testutil/keeper"
	"github.com/kiichain/kiichain3/x/payfi/keeper"
	"github.com/kiichain/kiichain3/x/payfi/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(t testing.TB) (*keeper.Keeper, types.MsgServer, context.Context) {
	k, ctx := keepertest.PayfiKeeper(t)
	return k, keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}

func setupMsgServerWithFundedAccounts(t testing.TB, accounts []sdk.AccAddress, amount sdk.Coins) (*keeper.Keeper, types.MsgServer, context.Context) {
	k, ctx := keepertest.PayfiKeeperWithFundedAccounts(t, accounts, amount)
	return k, keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}

func TestRegisterMerchant(t *testing.T) {
	_, msgServer, ctx := setupMsgServer(t)

	tests := []struct {
		name    string
		msg     *types.MsgRegisterMerchant
		wantErr bool
	}{
		{
			name: "valid registration",
			msg: &types.MsgRegisterMerchant{
				CreatorAddress:      "kii1x2w87cvt5mqjncav4lxy8yfreynn273xg764cm",
				BusinessName:        "Test Business",
				BusinessAddress:     "123 Test St",
				BusinessPhone:       "1234567890",
				BusinessEmail:       "test@test.com",
				BusinessDescription: "Test Description",
			},
			wantErr: false,
		},
		{
			name: "duplicate merchant",
			msg: &types.MsgRegisterMerchant{
				CreatorAddress:      "kii1x2w87cvt5mqjncav4lxy8yfreynn273xg764cm",
				BusinessName:        "Test Business",
				BusinessAddress:     "123 Test St",
				BusinessPhone:       "1234567890",
				BusinessEmail:       "test@test.com",
				BusinessDescription: "Test Description",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := msgServer.RegisterMerchant(ctx, tt.msg)
			if tt.wantErr {
				if err == nil {
					t.Error("RegisterMerchant() error = nil, wantErr = true")
				}
				return
			}
			if err != nil {
				t.Errorf("RegisterMerchant() error = %v, wantErr = false", err)
				return
			}
			if resp == nil {
				t.Error("RegisterMerchant() response is nil")
				return
			}
			if resp.MerchantAddress == "" {
				t.Error("RegisterMerchant() merchant address is empty")
			}
		})
	}
}

func TestVerifyMerchant(t *testing.T) {
	_, msgServer, ctx := setupMsgServer(t)

	// First register a merchant
	registerMsg := &types.MsgRegisterMerchant{
		CreatorAddress:      "kii1x2w87cvt5mqjncav4lxy8yfreynn273xg764cm",
		BusinessName:        "Test Business",
		BusinessAddress:     "123 Test St",
		BusinessPhone:       "1234567890",
		BusinessEmail:       "test@test.com",
		BusinessDescription: "Test Description",
	}
	registerResp, err := msgServer.RegisterMerchant(ctx, registerMsg)
	if err != nil {
		t.Fatalf("Failed to register merchant: %v", err)
	}

	tests := []struct {
		name    string
		msg     *types.MsgVerifyMerchant
		wantErr bool
	}{
		{
			name: "valid verification",
			msg: &types.MsgVerifyMerchant{
				KycAdminAddress: "kii1x2w87cvt5mqjncav4lxy8yfreynn273xg764cm",
				MerchantAddress: registerResp.MerchantAddress,
			},
			wantErr: false,
		},
		{
			name: "non-existent merchant",
			msg: &types.MsgVerifyMerchant{
				KycAdminAddress: "kii1x2w87cvt5mqjncav4lxy8yfreynn273xg764cm",
				MerchantAddress: "kii1nonexistent",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := msgServer.VerifyMerchant(ctx, tt.msg)
			if tt.wantErr {
				if err == nil {
					t.Error("VerifyMerchant() error = nil, wantErr = true")
				}
				return
			}
			if err != nil {
				t.Errorf("VerifyMerchant() error = %v, wantErr = false", err)
				return
			}
			if resp == nil {
				t.Error("VerifyMerchant() response is nil")
			}
		})
	}
}

func TestValidFlagChanged(t *testing.T) {
	k, msgServer, ctx := setupMsgServer(t)

	registerMsg := &types.MsgRegisterMerchant{
		CreatorAddress:      "kii1x2w87cvt5mqjncav4lxy8yfreynn273xg764cm",
		BusinessName:        "Test Business",
		BusinessAddress:     "123 Test St",
		BusinessPhone:       "1234567890",
		BusinessEmail:       "test@test.com",
		BusinessDescription: "Test Description",
	}
	registerResp, err := msgServer.RegisterMerchant(ctx, registerMsg)
	if err != nil {
		t.Fatalf("Failed to register merchant: %v", err)
	}

	require.False(t, k.GetMerchant(
		sdk.UnwrapSDKContext(ctx),
		sdk.MustAccAddressFromBech32(registerResp.MerchantAddress),
	).KycStatus)

	verifyMsg := &types.MsgVerifyMerchant{
		KycAdminAddress: "kii1x2w87cvt5mqjncav4lxy8yfreynn273xg764cm",
		MerchantAddress: registerResp.MerchantAddress,
	}
	_, err = msgServer.VerifyMerchant(ctx, verifyMsg)
	if err != nil {
		t.Fatalf("Failed to verify merchant: %v", err)
	}

	// Check if the merchant is valid
	merchant := k.GetMerchant(
		sdk.UnwrapSDKContext(ctx),
		sdk.MustAccAddressFromBech32(registerResp.MerchantAddress),
	)
	if merchant == nil {
		t.Fatalf("Merchant not found")
	}
	if !merchant.KycStatus {
		t.Fatalf("Merchant is not valid")
	}
}
func TestPayMerchantValidation(t *testing.T) {
	_, msgServer, ctx := setupMsgServerWithFundedAccounts(t, []sdk.AccAddress{
		sdk.MustAccAddressFromBech32("kii1x2w87cvt5mqjncav4lxy8yfreynn273xg764cm"),
	}, sdk.NewCoins(sdk.NewCoin("ukii", sdk.NewInt(1000000))))

	// First register a merchant
	registerMsg := &types.MsgRegisterMerchant{
		CreatorAddress:      "kii1x2w87cvt5mqjncav4lxy8yfreynn273xg764cm",
		BusinessName:        "Test Business",
		BusinessAddress:     "123 Test St",
		BusinessPhone:       "1234567890",
		BusinessEmail:       "test@test.com",
		BusinessDescription: "Test Description",
	}
	registerResp, err := msgServer.RegisterMerchant(ctx, registerMsg)
	if err != nil {
		t.Fatalf("Failed to register merchant: %v", err)
	}

	// Verify the merchant
	verifyMsg := &types.MsgVerifyMerchant{
		KycAdminAddress: "kii1x2w87cvt5mqjncav4lxy8yfreynn273xg764cm",
		MerchantAddress: registerResp.MerchantAddress,
	}
	_, err = msgServer.VerifyMerchant(ctx, verifyMsg)
	if err != nil {
		t.Fatalf("Failed to verify merchant: %v", err)
	}

	tests := []struct {
		name    string
		msg     *types.MsgPayMerchant
		wantErr bool
	}{
		{
			name: "valid payment message",
			msg: &types.MsgPayMerchant{
				PayerAddress:    "kii1x2w87cvt5mqjncav4lxy8yfreynn273xg764cm",
				MerchantAddress: registerResp.MerchantAddress,
				Amount:          sdk.NewCoins(sdk.NewCoin("ukii", sdk.NewInt(1000))),
			},
			wantErr: false,
		},
		{
			name: "invalid merchant address",
			msg: &types.MsgPayMerchant{
				PayerAddress:    "kii1x2w87cvt5mqjncav4lxy8yfreynn273xg764cm",
				MerchantAddress: "invalid_address",
				Amount:          sdk.NewCoins(sdk.NewCoin("ukii", sdk.NewInt(1000))),
			},
			wantErr: true,
		},
		{
			name: "zero amount",
			msg: &types.MsgPayMerchant{
				PayerAddress:    "kii1x2w87cvt5mqjncav4lxy8yfreynn273xg764cm",
				MerchantAddress: registerResp.MerchantAddress,
				Amount:          sdk.NewCoins(sdk.NewCoin("ukii", sdk.NewInt(0))),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := msgServer.PayMerchant(ctx, tt.msg)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}
		})
	}
}

func TestPayMerchantPaymentAndRevenue(t *testing.T) {
	k, msgServer, ctx := setupMsgServerWithFundedAccounts(t, []sdk.AccAddress{
		sdk.MustAccAddressFromBech32("kii1x2w87cvt5mqjncav4lxy8yfreynn273xg764cm"),
	}, sdk.NewCoins(sdk.NewCoin("ukii", sdk.NewInt(1000000))))

	// First register a merchant
	registerMsg := &types.MsgRegisterMerchant{
		CreatorAddress:      "kii1x2w87cvt5mqjncav4lxy8yfreynn273xg764cm",
		BusinessName:        "Test Business",
		BusinessAddress:     "123 Test St",
		BusinessPhone:       "1234567890",
		BusinessEmail:       "test@test.com",
		BusinessDescription: "Test Description",
	}
	registerResp, err := msgServer.RegisterMerchant(ctx, registerMsg)
	if err != nil {
		t.Fatalf("Failed to register merchant: %v", err)
	}

	// Verify the merchant
	verifyMsg := &types.MsgVerifyMerchant{
		KycAdminAddress: "kii1x2w87cvt5mqjncav4lxy8yfreynn273xg764cm",
		MerchantAddress: registerResp.MerchantAddress,
	}
	_, err = msgServer.VerifyMerchant(ctx, verifyMsg)
	if err != nil {
		t.Fatalf("Failed to verify merchant: %v", err)
	}

	// Make a payment
	paymentMsg := &types.MsgPayMerchant{
		PayerAddress:    "kii1x2w87cvt5mqjncav4lxy8yfreynn273xg764cm",
		MerchantAddress: registerResp.MerchantAddress,
		Amount:          sdk.NewCoins(sdk.NewCoin("ukii", sdk.NewInt(1000))),
	}

	resp, err := msgServer.PayMerchant(ctx, paymentMsg)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify payment was recorded
	merchantAddress := sdk.MustAccAddressFromBech32(paymentMsg.MerchantAddress)
	payment := k.GetPayment(sdk.UnwrapSDKContext(ctx), merchantAddress, 1)
	require.NotNil(t, payment)
	require.Equal(t, paymentMsg.PayerAddress, payment.PayerAddress)
	require.Equal(t, paymentMsg.MerchantAddress, payment.MerchantAddress)
	require.Equal(t, paymentMsg.Amount, payment.Amount)

	// Verify revenue was recorded
	revenue := k.GetMerchantRevenue(sdk.UnwrapSDKContext(ctx), merchantAddress)
	require.NotNil(t, revenue)
	require.Equal(t, paymentMsg.Amount, revenue.TotalRevenue)
	require.Equal(t, revenue.TotalPayments, int64(1))
}
