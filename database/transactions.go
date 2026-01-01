package database

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Mangrover007/banking-backend/internals/repository"
)

func (a *app) DepositMoney(user repository.User, amount int64, accountType repository.AccountType) error {
	account, err := a.query.FindAccount(a.ctx, repository.FindAccountParams{
		FkUserID: user.ID,
		Type:     accountType,
	})
	if err != nil {
		return err
	}

	// entry the transaction as pending
	trans, err := a.query.DepositTransaction(a.ctx, repository.DepositTransactionParams{
		Status:   repository.TransactionStatusPENDING,
		Amount:   amount,
		FkSender: account.ID,
	})
	if err != nil {
		return err
	}

	err = a.commitDeposit(trans.ID, account, amount)
	if err != nil {
		a.query.UpdateTransactionStatus(a.ctx, repository.UpdateTransactionStatusParams{
			ID:     trans.ID,
			Status: repository.TransactionStatusFAILED,
		})
		return err

		// what was that? what if this query itself fails?
		// at that point something has gone terrible wrong anyway so RIP i guess
	}

	return nil
}

func (a *app) commitDeposit(transactionID uuid.UUID, account repository.Account, amount int64) error {
	tx, err := a.conn.Begin(a.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(a.ctx)

	qtx := a.query.WithTx(tx)

	// update balance
	rows, err := qtx.UpdateBalanceDeposit(a.ctx, repository.UpdateBalanceDepositParams{
		ID:      account.ID,
		Balance: amount,
	})
	if err != nil {
		fmt.Println("Something went wrong oh no")
		return err
	}
	if rows != 1 {
		fmt.Println("Something different went wrong oh no oh god")
		return err
	}

	// update transaction status = SUCCESS
	rows, err = qtx.UpdateTransactionStatus(a.ctx, repository.UpdateTransactionStatusParams{
		ID:     transactionID,
		Status: repository.TransactionStatusSUCCESS,
	})
	if err != nil {
		fmt.Println("could not update transaction status")
		return err
	}
	if rows != 1 {
		fmt.Println("Could not update transaction status so FAIELD")
		return err // ! MAKE NEW ERROR
	}

	return tx.Commit(a.ctx)
}

func (a *app) WithdrawMoney(user repository.User, amount int64, accountType repository.AccountType) error {
	account, err := a.query.FindAccount(a.ctx, repository.FindAccountParams{
		FkUserID: user.ID,
		Type:     accountType,
	})
	if err != nil {
		return err
	}

	// entry the transaction as pending
	trans, err := a.query.WithdrawTransaction(a.ctx, repository.WithdrawTransactionParams{
		Status:   repository.TransactionStatusPENDING,
		Amount:   amount,
		FkSender: account.ID,
	})
	if err != nil {
		return err
	}

	err = a.commitWithdraw(trans.ID, account, amount)
	if err != nil {
		a.query.UpdateTransactionStatus(a.ctx, repository.UpdateTransactionStatusParams{
			ID:     trans.ID,
			Status: repository.TransactionStatusFAILED,
		})
		return err

		// what was that? what if this query itself fails?
		// at that point something has gone terrible wrong anyway so RIP i guess
	}

	return nil
}

func (a *app) commitWithdraw(transactionID uuid.UUID, account repository.Account, amount int64) error {
	tx, err := a.conn.Begin(a.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(a.ctx)

	qtx := a.query.WithTx(tx)

	// update balance
	rows, err := qtx.UpdateBalanceWithdraw(a.ctx, repository.UpdateBalanceWithdrawParams{
		ID:      account.ID,
		Balance: amount,
	})
	if err != nil {
		fmt.Println("could not update balance")
		return err
	}
	if rows != 1 {
		fmt.Println("could not update balance again")
		return err
	}

	// update the transaction status to SUCCESS
	rows, err = qtx.UpdateTransactionStatus(a.ctx, repository.UpdateTransactionStatusParams{
		ID:     transactionID,
		Status: repository.TransactionStatusSUCCESS,
	})
	if err != nil {
		fmt.Println("could not update transaction status")
		return err
	}
	if rows != 1 {
		fmt.Println("Could not update transaction status so FAIELD")
		return err // ! MAKE NEW ERROR
	}

	return tx.Commit(a.ctx)
}

func (a *app) TransferMoney(sender, recipient repository.Account, amount int64) error {

	trans, err := a.query.TransferTransaction(a.ctx, repository.TransferTransactionParams{
		Status:      repository.TransactionStatusPENDING,
		Amount:      amount,
		FkSender:    sender.ID,
		FkRecipient: pgtype.UUID{Bytes: recipient.ID, Valid: true},
	})

	err = a.commitTransfer(sender, recipient, amount, trans)
	if err != nil {
		a.query.UpdateTransactionStatus(a.ctx, repository.UpdateTransactionStatusParams{
			ID:     trans.ID,
			Status: repository.TransactionStatusFAILED,
		})
		return err
	}

	return nil
}

func (a *app) commitTransfer(sender, recipient repository.Account, amount int64, trans repository.Transaction) error {
	tx, err := a.conn.Begin(a.ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(a.ctx)

	qtx := a.query.WithTx(tx)

	// obtain locks on SENDER and RECIPIENT accounts
	// prevent deadlock by obtaining the locks in consistent order
	if sender.ID.String() > recipient.ID.String() {
		lock_s, err := qtx.LockTransferAccount(a.ctx, sender.ID)
		if err != nil {
			return err
		}
		if lock_s != 1 {
			return fmt.Errorf("Could NOT obtain XLOCK on SENDER")
		}
		lock_r, err := qtx.LockTransferAccount(a.ctx, recipient.ID)
		if err != nil {
			return err
		}
		if lock_r != 1 {
			return fmt.Errorf("Could NOT obtain XLOCK on RECIPIENT")
		}
	} else {
		lock_r, err := qtx.LockTransferAccount(a.ctx, recipient.ID)
		if err != nil {
			return err
		}
		if lock_r != 1 {
			return fmt.Errorf("Could NOT obtain XLOCK on RECIPIENT")
		}
		lock_s, err := qtx.LockTransferAccount(a.ctx, sender.ID)
		if err != nil {
			return err
		}
		if lock_s != 1 {
			return fmt.Errorf("Could NOT obtain XLOCK on SENDER")
		}
	}

	// deduct money from sender
	rows_s, err := qtx.UpdateBalanceWithdraw(a.ctx, repository.UpdateBalanceWithdrawParams{
		ID:      sender.ID,
		Balance: amount,
	})
	if err != nil {
		return err
	}
	if rows_s != 1 {
		return fmt.Errorf("Could NOT UPDATE BALANCE for SENDER")
	}

	// add moeny to receiver
	rows_r, err := qtx.UpdateBalanceDeposit(a.ctx, repository.UpdateBalanceDepositParams{
		ID:      recipient.ID,
		Balance: amount,
	})
	if err != nil {
		return err
	}
	if rows_r != 1 {
		return fmt.Errorf("Could NOT UPDATE BALANCE for RECIPIENT")
	}

	qtx.UpdateTransactionStatus(a.ctx, repository.UpdateTransactionStatusParams{
		ID:     trans.ID,
		Status: repository.TransactionStatusSUCCESS,
	})

	// commit transaction
	return tx.Commit(a.ctx)
}
