package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testkeeper "github.com/kiichain/kiichain3/testutil/keeper"
	"github.com/kiichain/kiichain3/x/payfi/types"
	"github.com/stretchr/testify/require"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := testkeeper.PayfiKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	response, err := keeper.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}

func TestAllPaymentsQuery(t *testing.T) {
	keeper, ctx := testkeeper.PayfiKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)

	// Create test payments
	testMerchant := sdk.MustAccAddressFromBech32("kii1x2w87cvt5mqjncav4lxy8yfreynn273xg764cm")
	payment1 := types.Payment{
		PaymentId:       1,
		MerchantAddress: testMerchant.String(),
		PayerAddress:    "kii1payer1",
		Amount:          sdk.NewCoins(sdk.NewCoin("ukii", sdk.NewInt(1000))),
	}
	payment2 := types.Payment{
		PaymentId:       2,
		MerchantAddress: testMerchant.String(),
		PayerAddress:    "kii1payer2",
		Amount:          sdk.NewCoins(sdk.NewCoin("ukii", sdk.NewInt(2000))),
	}

	keeper.SetPayment(ctx, payment1)
	keeper.SetPayment(ctx, payment2)

	response, err := keeper.AllPayments(wctx, &types.QueryAllPaymentsRequest{})
	require.NoError(t, err)
	require.Len(t, response.Payments, 2)
}

func TestPaymentQuery(t *testing.T) {
	keeper, ctx := testkeeper.PayfiKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)

	testMerchant := sdk.MustAccAddressFromBech32("kii1x2w87cvt5mqjncav4lxy8yfreynn273xg764cm")
	payment := types.Payment{
		PaymentId:       1,
		MerchantAddress: testMerchant.String(),
		PayerAddress:    "kii1payer1",
		Amount:          sdk.NewCoins(sdk.NewCoin("ukii", sdk.NewInt(1000))),
	}
	keeper.SetPayment(ctx, payment)

	tests := []struct {
		name      string
		request   *types.QueryPaymentRequest
		expected  *types.Payment
		shouldErr bool
	}{
		{
			name: "Valid Payment Query",
			request: &types.QueryPaymentRequest{
				PaymentId:       1,
				MerchantAddress: testMerchant.String(),
			},
			expected:  &payment,
			shouldErr: false,
		},
		{
			name: "Non-existent Payment",
			request: &types.QueryPaymentRequest{
				PaymentId:       999,
				MerchantAddress: testMerchant.String(),
			},
			expected:  nil,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := keeper.Payment(wctx, tt.request)
			if tt.shouldErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expected, response.Payment)
		})
	}
}

func TestMerchantQuery(t *testing.T) {
	keeper, ctx := testkeeper.PayfiKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)

	testMerchant := sdk.MustAccAddressFromBech32("kii1x2w87cvt5mqjncav4lxy8yfreynn273xg764cm")
	merchant := types.Merchant{
		MerchantAddress: testMerchant.String(),
		BusinessName:    "Test Business",
		BusinessAddress: "123 Test St",
		BusinessPhone:   "1234567890",
		BusinessEmail:   "test@test.com",
		KycStatus:       true,
	}
	keeper.SetMerchant(ctx, merchant)

	response, err := keeper.Merchant(wctx, &types.QueryMerchantRequest{
		MerchantAddress: testMerchant.String(),
	})
	require.NoError(t, err)
	require.Equal(t, &merchant, response.Merchant)
}

func TestAllMerchantsQuery(t *testing.T) {
	keeper, ctx := testkeeper.PayfiKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)

	// Create test merchants
	testMerchant1 := testkeeper.RandomAccAddress()
	testMerchant2 := testkeeper.RandomAccAddress()

	merchant1 := types.Merchant{
		MerchantAddress: testMerchant1.String(),
		BusinessName:    "Test Business 1",
		KycStatus:       true,
	}
	merchant2 := types.Merchant{
		MerchantAddress: testMerchant2.String(),
		BusinessName:    "Test Business 2",
		KycStatus:       true,
	}

	keeper.SetMerchant(ctx, merchant1)
	keeper.SetMerchant(ctx, merchant2)

	response, err := keeper.AllMerchants(wctx, &types.QueryAllMerchantsRequest{})
	require.NoError(t, err)
	require.Len(t, response.Merchants, 2)
}

func TestMerchantRevenueQuery(t *testing.T) {
	keeper, ctx := testkeeper.PayfiKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)

	testMerchant := sdk.MustAccAddressFromBech32("kii1x2w87cvt5mqjncav4lxy8yfreynn273xg764cm")
	revenue := types.MerchantRevenue{
		MerchantAddress: testMerchant.String(),
		TotalRevenue:    sdk.NewCoins(sdk.NewCoin("ukii", sdk.NewInt(5000))),
		TotalPayments:   5,
	}
	keeper.SetMerchantRevenue(ctx, revenue)

	response, err := keeper.MerchantRevenue(wctx, &types.QueryMerchantRevenueRequest{
		MerchantAddress: testMerchant.String(),
	})
	require.NoError(t, err)
	require.Equal(t, &revenue, response.Revenue)
}

func TestAllRevenueQuery(t *testing.T) {
	keeper, ctx := testkeeper.PayfiKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)

	// Create test revenues
	testMerchant1 := testkeeper.RandomAccAddress()
	testMerchant2 := testkeeper.RandomAccAddress()

	revenue1 := types.MerchantRevenue{
		MerchantAddress: testMerchant1.String(),
		TotalRevenue:    sdk.NewCoins(sdk.NewCoin("ukii", sdk.NewInt(5000))),
		TotalPayments:   5,
	}
	revenue2 := types.MerchantRevenue{
		MerchantAddress: testMerchant2.String(),
		TotalRevenue:    sdk.NewCoins(sdk.NewCoin("ukii", sdk.NewInt(3000))),
		TotalPayments:   3,
	}

	keeper.SetMerchantRevenue(ctx, revenue1)
	keeper.SetMerchantRevenue(ctx, revenue2)

	response, err := keeper.AllRevenue(wctx, &types.QueryAllRevenueRequest{})
	require.NoError(t, err)
	require.Len(t, response.Revenues, 2)
}
