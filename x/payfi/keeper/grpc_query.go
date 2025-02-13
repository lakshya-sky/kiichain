package keeper

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kiichain/kiichain3/x/payfi/types"
)

var _ types.QueryServer = Keeper{}

// Params implements the Query/Params gRPC method
func (k Keeper) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	return &types.QueryParamsResponse{Params: k.GetParams(sdkCtx)}, nil
}

// AllPayments implements the Query/AllPayments gRPC method
func (k Keeper) AllPayments(ctx context.Context, req *types.QueryAllPaymentsRequest) (*types.QueryAllPaymentsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(k.storeKey)
	paymentStore := prefix.NewStore(store, types.KeyPrefix(types.PaymentKey))

	var payments []*types.Payment
	pageRes, err := query.Paginate(paymentStore, req.Pagination, func(key []byte, value []byte) error {
		var payment *types.Payment
		if err := k.cdc.Unmarshal(value, payment); err != nil {
			return err
		}
		payments = append(payments, payment)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPaymentsResponse{
		Payments:   payments,
		Pagination: pageRes,
	}, nil
}

// Payment implements the Query/Payment gRPC method
func (k Keeper) Payment(ctx context.Context, req *types.QueryPaymentRequest) (*types.QueryPaymentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	merchantAddress, err := sdk.AccAddressFromBech32(req.MerchantAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid merchant address")
	}
	paymentId, err := strconv.ParseInt(req.PaymentId, 10, 64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid payment id")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	payment := k.GetPayment(sdkCtx, merchantAddress, paymentId)
	if payment == nil {
		return nil, status.Error(codes.NotFound, "payment not found")
	}

	return &types.QueryPaymentResponse{Payment: payment}, nil
}

// MerchantPayments implements the Query/MerchantPayments gRPC method
func (k Keeper) MerchantPayments(ctx context.Context, req *types.QueryMerchantPaymentsRequest) (*types.QueryMerchantPaymentsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	merchantAddress, err := sdk.AccAddressFromBech32(req.MerchantAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid merchant address")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(k.storeKey)
	merchantPaymentStore := prefix.NewStore(store, types.MerchantKeyPrefix(merchantAddress))

	var payments []*types.Payment
	pageRes, err := query.Paginate(merchantPaymentStore, req.Pagination, func(key []byte, value []byte) error {
		var payment *types.Payment
		if err := k.cdc.Unmarshal(value, payment); err != nil {
			return err
		}
		payments = append(payments, payment)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryMerchantPaymentsResponse{
		Payments:   payments,
		Pagination: pageRes,
	}, nil
}

// MerchantRevenue implements the Query/MerchantRevenue gRPC method
func (k Keeper) MerchantRevenue(ctx context.Context, req *types.QueryMerchantRevenueRequest) (*types.QueryMerchantRevenueResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	merchantAddress, err := sdk.AccAddressFromBech32(req.MerchantAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid merchant address")
	}
	revenue := k.GetMerchantRevenue(sdkCtx, merchantAddress)
	if revenue == nil {
		return nil, status.Error(codes.NotFound, "merchant revenue not found")
	}

	return &types.QueryMerchantRevenueResponse{Revenue: revenue}, nil
}

// Merchant implements the Query/Merchant gRPC method
func (k Keeper) Merchant(ctx context.Context, req *types.QueryMerchantRequest) (*types.QueryMerchantResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	merchantAddress, err := sdk.AccAddressFromBech32(req.MerchantAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid merchant address")
	}
	merchant := k.GetMerchant(sdkCtx, merchantAddress)
	if merchant == nil {
		return nil, status.Error(codes.NotFound, "merchant not found")
	}

	return &types.QueryMerchantResponse{Merchant: merchant}, nil
}

// AllMerchants implements the Query/AllMerchants gRPC method
func (k Keeper) AllMerchants(ctx context.Context, req *types.QueryAllMerchantsRequest) (*types.QueryAllMerchantsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(k.storeKey)
	merchantStore := prefix.NewStore(store, types.KeyPrefix(types.MerchantKey))

	var merchants []*types.Merchant
	pageRes, err := query.Paginate(merchantStore, req.Pagination, func(key []byte, value []byte) error {
		var merchant *types.Merchant
		if err := k.cdc.Unmarshal(value, merchant); err != nil {
			return err
		}
		merchants = append(merchants, merchant)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllMerchantsResponse{
		Merchants:  merchants,
		Pagination: pageRes,
	}, nil
}

// AllRevenue implements the Query/AllRevenue gRPC method
func (k Keeper) AllRevenue(ctx context.Context, req *types.QueryAllRevenueRequest) (*types.QueryAllRevenueResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(k.storeKey)
	revenueStore := prefix.NewStore(store, types.KeyPrefix(types.MerchantRevenueKey))

	var revenues []*types.MerchantRevenue
	pageRes, err := query.Paginate(revenueStore, req.Pagination, func(key []byte, value []byte) error {
		var revenue *types.MerchantRevenue
		if err := k.cdc.Unmarshal(value, revenue); err != nil {
			return err
		}
		revenues = append(revenues, revenue)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllRevenueResponse{
		Revenues:   revenues,
		Pagination: pageRes,
	}, nil
}
