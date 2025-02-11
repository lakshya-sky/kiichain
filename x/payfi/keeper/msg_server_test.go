package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/kiichain/kiichain3/testutil/keeper"
	"github.com/kiichain/kiichain3/x/payfi/keeper"
	"github.com/kiichain/kiichain3/x/payfi/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.PayfiKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
