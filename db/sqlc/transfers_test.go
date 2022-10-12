package db

import (
	"context"
	"simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, acc1, acc2 Account) Transfer {
	args := CreateTransferParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        int64(util.RandomAmount()),
	}

	trans, err := testQueries.CreateTransfer(context.Background(), args)

	require.NoError(t, err)
	require.Equal(t, args.FromAccountID, trans.FromAccountID)
	require.Equal(t, args.ToAccountID, trans.ToAccountID)
	require.Equal(t, args.Amount, trans.Amount)

	require.NotZero(t, trans.ID)
	require.NotZero(t, trans.CreatedAt)

	return trans
}

func TestCreatTransfer(t *testing.T) {
	acc1 := creatRandomAccount(t)
	acc2 := creatRandomAccount(t)
	createRandomTransfer(t, acc1, acc2)
}

func TestGetTransfer(t *testing.T) {
	acc1 := creatRandomAccount(t)
	acc2 := creatRandomAccount(t)
	trans1 := createRandomTransfer(t, acc1, acc2)
	trans2, err := testQueries.GetTransfer(context.Background(), trans1.ID)

	require.NoError(t, err)
	require.Equal(t, trans1.ID, trans2.ID)
	require.Equal(t, trans1.FromAccountID, trans2.FromAccountID)
	require.Equal(t, trans1.ToAccountID, trans2.ToAccountID)
	require.Equal(t, trans1.Amount, trans2.Amount)

	require.WithinDuration(t, trans1.CreatedAt, trans2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	acc1 := creatRandomAccount(t)
	acc2 := creatRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomTransfer(t, acc1, acc2)
	}

	args := ListTransfersParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, trans := range transfers {
		require.NotEmpty(t, trans)
		require.Equal(t, acc1.ID, trans.FromAccountID)
		require.Equal(t, acc2.ID, trans.ToAccountID)
	}
}
