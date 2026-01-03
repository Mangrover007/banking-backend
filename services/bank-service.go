package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Mangrover007/banking-backend/internals/repository"
)

type BankService interface {
	Deposit(ctx context.Context, account repository.Account, amount int64) error
	Withdraw(ctx context.Context, account repository.Account, amount int64) error
	Transfer(ctx context.Context, sender, recipient repository.Account, amount int64) error
	OpenSavingsAccount(ctx context.Context, user repository.User, balance int64, _type repository.AccountType) (repository.Account, error)
	OpenCheckingAccount(ctx context.Context, user repository.User, balance int64, _type repository.AccountType) (repository.Account, error)
	FindAccountByID(ctx context.Context, accountID uuid.UUID) (repository.Account, error)
	FindAccount(ctx context.Context, userID uuid.UUID, _type repository.AccountType) (repository.Account, error)
}

type bank struct {
	conn  *pgx.Conn
	query *repository.Queries
}

var (
	ErrRowsCorruption = errors.New("More than 1 or no rows were affected during the transaction. Rolling back.")
)

func NewBank(conn *pgx.Conn) BankService {
	return &bank{
		conn:  conn,
		query: repository.New(conn),
	}
}

func (b *bank) Deposit(ctx context.Context, account repository.Account, amount int64) error {
	// create transaction record
	trans, err := b.query.DepositTransaction(ctx, repository.DepositTransactionParams{
		Status:   repository.TransactionStatusPENDING,
		Amount:   amount,
		FkSender: account.ID,
	})
	if err != nil {
		return err
	}

	err = func() error {
		tx, err := b.conn.BeginTx(ctx, pgx.TxOptions{
			IsoLevel:   pgx.Serializable,
			AccessMode: pgx.ReadWrite,
		})
		if err != nil {
			return err
		}
		defer tx.Rollback(ctx)

		qtx := b.query.WithTx(tx)

		// attempt deposit
		rowsAff, err := b.depositWrapper(ctx, qtx, amount, account.ID)
		if err != nil {
			tx.Rollback(ctx) // paranoia rollback
			return err
		}
		if rowsAff != 1 {
			tx.Rollback(ctx)         // another paranoia rollback
			return ErrRowsCorruption // how does returning this error help
		}

		// if commit fails, then no changes were made duh!
		err = tx.Commit(ctx)
		return err
	}()

	if err != nil {
		b.updateTransactionWrapper(ctx, trans.ID, repository.TransactionStatusFAILED)
	} else {
		b.updateTransactionWrapper(ctx, trans.ID, repository.TransactionStatusSUCCESS)
	}

	return err
}

func (b *bank) Withdraw(ctx context.Context, account repository.Account, amount int64) error {
	// insert transaction record
	trans, err := b.query.WithdrawTransaction(ctx, repository.WithdrawTransactionParams{
		Status:   repository.TransactionStatusPENDING,
		Amount:   amount,
		FkSender: account.ID,
	})
	if err != nil {
		return err
	}

	err = func() error {
		tx, err := b.conn.BeginTx(ctx, pgx.TxOptions{
			IsoLevel:   pgx.Serializable,
			AccessMode: pgx.ReadWrite,
		})
		if err != nil {
			return err
		}
		defer tx.Rollback(ctx)

		qtx := b.query.WithTx(tx)

		// try withdraw
		rowsAff, err := b.withdrawWrapper(ctx, qtx, amount, account.ID)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
		if rowsAff != 1 {
			tx.Rollback(ctx)
			return ErrRowsCorruption
		}

		return tx.Commit(ctx)
	}()

	if err != nil {
		b.updateTransactionWrapper(ctx, trans.ID, repository.TransactionStatusFAILED)
	} else {
		b.updateTransactionWrapper(ctx, trans.ID, repository.TransactionStatusSUCCESS)
	}

	return err
}

func (b *bank) Transfer(ctx context.Context, sender, recipient repository.Account, amount int64) error {
	// insert transfer transaction record
	trans, err := b.query.TransferTransaction(ctx, repository.TransferTransactionParams{
		FkSender:    sender.ID,
		FkRecipient: pgtype.UUID{Bytes: recipient.ID, Valid: true},
		Amount:      amount,
	})
	if err != nil {
		return err
	}

	err = func() error {
		tx, err := b.conn.BeginTx(ctx, pgx.TxOptions{
			IsoLevel:   pgx.Serializable,
			AccessMode: pgx.ReadWrite,
		})
		if err != nil {
			return err
		}

		qtx := b.query.WithTx(tx)

		// lock accounts in the same order
		// otherwise may deadlock
		b.lockAccountsForTransfer(ctx, qtx, sender, recipient)

		// attempt transfer
		rowsAff, err := b.withdrawWrapper(ctx, qtx, amount, sender.ID)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
		if rowsAff != 1 {
			tx.Rollback(ctx)
			return ErrRowsCorruption
		}

		rowsAff, err = b.depositWrapper(ctx, qtx, amount, recipient.ID)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
		if rowsAff != 1 {
			tx.Rollback(ctx)
			return ErrRowsCorruption
		}

		return tx.Commit(ctx)
	}()

	if err != nil {
		b.updateTransactionWrapper(ctx, trans.ID, repository.TransactionStatusFAILED)
	} else {
		b.updateTransactionWrapper(ctx, trans.ID, repository.TransactionStatusSUCCESS)
	}

	return err
}

func (b *bank) OpenSavingsAccount(ctx context.Context, user repository.User, balance int64, _type repository.AccountType) (repository.Account, error) {
	tx, err := b.conn.Begin(ctx)
	if err != nil {
		return repository.Account{}, err
	}
	qtx := b.query.WithTx(tx)

	account, err := qtx.OpenAccount(ctx, repository.OpenAccountParams{
		Balance:  balance,
		Type:     repository.AccountTypeCHECKING,
		FkUserID: user.ID,
	})

	var interestRate pgtype.Numeric
	interestRate.Scan("0.4")
	_, err = qtx.OpenSavingsAccount(ctx, repository.OpenSavingsAccountParams{
		AccountID:       account.ID,
		InterestRate:    interestRate,
		MinBalance:      25,
		WithdrawalLimit: 5000,
	})
	if err != nil {
		tx.Rollback(ctx)
		return repository.Account{}, err
	}

	return account, tx.Commit(ctx)
}

func (b *bank) OpenCheckingAccount(ctx context.Context, user repository.User, balance int64, _type repository.AccountType) (repository.Account, error) {
	tx, err := b.conn.Begin(ctx)
	if err != nil {
		return repository.Account{}, err
	}
	qtx := b.query.WithTx(tx)

	account, err := qtx.OpenAccount(ctx, repository.OpenAccountParams{
		Balance:  balance,
		Type:     repository.AccountTypeSAVINGS,
		FkUserID: user.ID,
	})

	var overdraftLimit pgtype.Numeric
	overdraftLimit.Scan("0.4")
	_, err = qtx.OpenCheckingAccount(ctx, repository.OpenCheckingAccountParams{
		AccountID:      account.ID,
		OverdraftLimit: overdraftLimit,
		MaintenanceFee: 12,
	})
	if err != nil {
		tx.Rollback(ctx)
		return repository.Account{}, err
	}

	return account, tx.Commit(ctx)
}

func Shit() { fmt.Println("eat shit") }

// HELPER FUNCTIONS
func (b *bank) depositWrapper(ctx context.Context, qtx *repository.Queries, amount int64, id uuid.UUID) (int64, error) {
	rowsAff, err := qtx.UpdateBalanceDeposit(ctx, repository.UpdateBalanceDepositParams{
		ID:      id,
		Balance: amount,
	})
	return rowsAff, err
}

func (b *bank) withdrawWrapper(ctx context.Context, qtx *repository.Queries, amount int64, id uuid.UUID) (int64, error) {
	rowsAff, err := qtx.UpdateBalanceWithdraw(ctx, repository.UpdateBalanceWithdrawParams{
		ID:      id,
		Balance: amount,
	})
	return rowsAff, err
}

func (b *bank) updateTransactionWrapper(ctx context.Context, transID uuid.UUID, status repository.TransactionStatus) (int64, error) {
	rowsAff, err := b.query.UpdateTransactionStatus(ctx, repository.UpdateTransactionStatusParams{
		Status: status,
		ID:     transID,
	})
	return rowsAff, err
}

func (b *bank) lockAccountsForTransfer(ctx context.Context, qtx *repository.Queries, sender, recipient repository.Account) error {
	if sender.ID.String() > recipient.ID.String() {
		lock, err := qtx.LockTransferAccount(ctx, sender.ID)
		if err != nil {
			return err
		}
		if lock != 1 {
			return ErrRowsCorruption
		}
		lock, err = qtx.LockTransferAccount(ctx, recipient.ID)
		if err != nil {
			return err
		}
		if lock != 1 {
			return ErrRowsCorruption
		}
	} else {
		lock, err := qtx.LockTransferAccount(ctx, recipient.ID)
		if err != nil {
			return err
		}
		if lock != 1 {
			return ErrRowsCorruption
		}
		lock, err = qtx.LockTransferAccount(ctx, sender.ID)
		if err != nil {
			return err
		}
		if lock != 1 {
			return ErrRowsCorruption
		}
	}
	return nil
}

func (b *bank) FindAccount(ctx context.Context, userID uuid.UUID, _type repository.AccountType) (repository.Account, error) {
	account, err := b.query.FindAccount(ctx, repository.FindAccountParams{
		FkUserID: userID,
		Type:     _type,
	})
	if err != nil {
		return repository.Account{}, err
	}

	return account, nil
}
func (b *bank) FindAccountByID(ctx context.Context, accountID uuid.UUID) (repository.Account, error) {
	account, err := b.query.FindAccountByID(ctx, accountID)
	if err != nil {
		return repository.Account{}, err
	}
	return account, nil
}
