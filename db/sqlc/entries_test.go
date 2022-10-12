package db

import (
	"context"
	"simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, acc Account) Entry {
	args := CreateEntryParams{
		AccountID: acc.ID,
		Amount:    int64(util.RandomAmount()),
	}

	ent, err := testQueries.CreateEntry(context.Background(), args)

	require.NoError(t, err)
	require.Equal(t, args.AccountID, ent.AccountID)
	require.Equal(t, args.Amount, ent.Amount)

	require.NotZero(t, ent.ID)
	require.NotZero(t, ent.CreatedAt)

	return ent
}

func TestCreatEntry(t *testing.T) {
	acc := creatRandomAccount(t)
	createRandomEntry(t, acc)
}

func TestGetEntry(t *testing.T) {
	acc := creatRandomAccount(t)
	ent1 := createRandomEntry(t, acc)
	ent2, err := testQueries.GetEntry(context.Background(), ent1.ID)

	require.NoError(t, err)
	require.Equal(t, ent1.ID, ent2.ID)
	require.Equal(t, ent1.AccountID, ent2.AccountID)
	require.Equal(t, ent1.Amount, ent2.Amount)

	require.WithinDuration(t, ent1.CreatedAt, ent2.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	acc := creatRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, acc)
	}

	args := ListEntriesParams{
		AccountID: acc.ID,
		Limit:     5,
		Offset:    5,
	}

	ents, err := testQueries.ListEntries(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, ents, 5)

	for _, e := range ents {
		require.NotEmpty(t, e)
		require.Equal(t, acc.ID, e.AccountID)
	}
}
