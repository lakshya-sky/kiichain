package keeper

import (
	"github.com/kiichain/kiichain3/x/payfi/types"
)

var _ types.QueryServer = Keeper{}
