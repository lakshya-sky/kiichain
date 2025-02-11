package payfi_test

import (
	"testing"

	keepertest "github.com/kiichain/kiichain3/testutil/keeper"
	"github.com/kiichain/kiichain3/testutil/nullify"
	"github.com/kiichain/kiichain3/x/payfi"
	"github.com/kiichain/kiichain3/x/payfi/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.PayfiKeeper(t)
	payfi.InitGenesis(ctx, *k, genesisState)
	got := payfi.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
