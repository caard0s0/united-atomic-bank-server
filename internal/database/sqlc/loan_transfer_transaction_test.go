package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoanTransferTransaction(t *testing.T) {
	store := NewStore(testDB)

	account := createRandomAccount(t)
	fmt.Println(">> before", account.Balance)

	// run a concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan LoanTransferTransactionResult)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := store.LoanTransferTransaction(ctx, CreateLoanTransferParams{
				AccountID: account.ID,
				Amount:    amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer loan
		loan := result.Loan
		require.NotEmpty(t, loan)
		require.Equal(t, account.ID, loan.AccountID)
		require.Equal(t, amount, loan.Amount)
		require.NotZero(t, loan.ID)
		require.NotZero(t, loan.StartAt)
		require.NotZero(t, loan.EndAt)

		_, err = store.GetLoanTransfer(context.Background(), loan.ID)
		require.NoError(t, err)

		// check entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check account
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account.ID, toAccount.ID)
	}

	// check the final updated balance
	updatedAccount, err := TestQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)

	fmt.Println(">> after", updatedAccount.Balance)
	require.Equal(t, account.Balance+int64(n)*amount, updatedAccount.Balance)
}
