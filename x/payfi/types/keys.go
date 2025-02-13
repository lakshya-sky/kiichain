package types

import (
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
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func MerchantKeyPrefix(merchantAddress sdk.AccAddress) []byte {
	return append(KeyPrefix(MerchantKey), merchantAddress.Bytes()...)
}
