package db

import (
	"context"
	"database/sql"
	"simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func creatRandomAccount(t *testing.T) Account {
	user := creatRandomUser(t)
	args := CreateAccountParams{
		Owner:    user.Username,
		Balance:  int64(util.RandomAmount()),
		Currency: util.RandomCurrency(),
	}

	acc, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.Equal(t, args.Owner, acc.Owner)
	require.Equal(t, args.Balance, acc.Balance)
	require.Equal(t, args.Currency, acc.Currency)

	require.NotZero(t, acc.ID)
	require.NotZero(t, acc.CreatedAt)

	return acc
}

func TestCreatAccount(t *testing.T) {
	creatRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc1 := creatRandomAccount(t)
	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)

	require.NoError(t, err)
	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, acc1.Balance, acc2.Balance)
	require.Equal(t, acc1.Currency, acc2.Currency)

	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	acc1 := creatRandomAccount(t)
	args := UpdateAccountParams{
		ID:      acc1.ID,
		Balance: int64(util.RandomAmount()),
	}

	acc2, err := testQueries.UpdateAccount(context.Background(), args)

	require.NoError(t, err)
	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, args.Balance, acc2.Balance)
	require.Equal(t, acc1.Currency, acc2.Currency)

	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)
}

func TestAddAccountBalance(t *testing.T) {
	acc1 := creatRandomAccount(t)
	args := AddAccountBalanceParams{
		ID:     acc1.ID,
		Amount: int64(util.RandomAmount()),
	}

	acc2, err := testQueries.AddAccountBalance(context.Background(), args)

	require.NoError(t, err)
	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, acc2.Balance, acc1.Balance+args.Amount)
	require.Equal(t, acc1.Currency, acc2.Currency)

	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	acc1 := creatRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	acc2, getErr := testQueries.GetAccount(context.Background(), acc1.ID)
	require.Error(t, getErr)
	require.EqualError(t, getErr, sql.ErrNoRows.Error())
	require.Empty(t, acc2)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = creatRandomAccount(t)
	}

	args := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accs, err := testQueries.ListAccounts(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, accs)

	for _, a := range accs {
		require.NotEmpty(t, a)
		require.Equal(t, lastAccount.Owner, a.Owner)
	}
}
