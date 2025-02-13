package types

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "payfi"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_payfi"

	// MerchantKey defines the key for the merchant
	MerchantKey = "merchant"

	// PaymentKey defines the key for the payment
	PaymentKey = "payment"

	// MerchantRevenueKey defines the key for the merchant revenue
	MerchantRevenueKey = "merchant_revenue"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func MerchantKeyPrefix(merchantAddress sdk.AccAddress) []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, KeyPrefix(MerchantKey)...)
	bytes = append(bytes, merchantAddress.Bytes()...)
	return bytes
}

func PaymentKeyPrefix(merchantAddress sdk.AccAddress, paymentId int64) []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, KeyPrefix(PaymentKey)...)
	bytes = append(bytes, merchantAddress.Bytes()...)
	bytes = append(bytes, []byte(strconv.FormatInt(paymentId, 10))...)
	return bytes
}

func MerchantRevenueKeyPrefix(merchantAddress sdk.AccAddress) []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, KeyPrefix(MerchantRevenueKey)...)
	bytes = append(bytes, merchantAddress.Bytes()...)
	return bytes
}
