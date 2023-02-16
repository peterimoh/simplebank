package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/peterimoh/utils"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) Entry {

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    utils.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry1 := createRandomEntry(t, account)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.AccountID)
	require.NoError(t, err)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)

}

func TestListEntries(t *testing.T) {
	for i := 0; i < 6; i++ {
		createRandomEntry(t, createRandomAccount(t))
	}

	arg := ListEntriesParams{
		Limit:  3,
		Offset: 2,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 3)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

func TestDeleteEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry1 := createRandomEntry(t, account)

	err := testQueries.DeleteEntry(context.Background(), entry1.ID)
	require.NoError(t, err)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry2)
}

func TestUpdateEntry(t *testing.T) {
	account := createRandomAccount(t)

	entry1 := createRandomEntry(t, account)

	args := UpdateEntryParams{
		ID:     entry1.ID,
		Amount: utils.RandomMoney(),
	}

	err := testQueries.UpdateEntry(context.Background(), args)
	require.NoError(t, err)
}
