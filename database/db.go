package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"

	"github.com/Mangrover007/banking-backend/internals/repository"
)

type app struct {
	conn  *pgx.Conn
	ctx   context.Context
	query *repository.Queries
}

var (
	App *app
)

func NewApp() (*app, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	DB_URL := os.Getenv("DATABASE_URL")

	conn, err := pgx.Connect(ctx, DB_URL)
	if err != nil {
		fmt.Println("Could not connect to DB")
		return nil, err
	}

	query := repository.New(conn)

	// make sure the connection is working
	users, _ := query.FindAllUsers(ctx)
	for i, user := range users {
		fmt.Printf("%d: %+v\n", i, user)
	}

	App = &app{
		conn:  conn,
		ctx:   ctx,
		query: query,
	}

	return App, nil
}

func (a *app) RegisterUser(firstName, lastName, email, phoneNumber, address, password string) (repository.User, error) {
	user, err := a.query.CreateUser(a.ctx, repository.CreateUserParams{
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		PhoneNumber: phoneNumber,
		Address:     address,
		Password:    password, // encrypt this pls
	})

	if err != nil {
		return repository.User{}, err
	}

	return user, nil
}

func (a *app) CreateAccount(user repository.User, initialBalance int64, accountType repository.AccountType) (repository.Account, error) {
	tx, err := a.conn.Begin(a.ctx)
	if err != nil {
		tx.Rollback(a.ctx)
		return repository.Account{}, err
	}

	defer tx.Rollback(a.ctx)

	qtx := a.query.WithTx(tx)

	account, err := qtx.OpenAccount(a.ctx, repository.OpenAccountParams{
		AccountNumber: "109876543210",
		Balance:       initialBalance,
		Type:          accountType,
		FkUserID:      user.ID,
	})

	if err != nil {
		return repository.Account{}, err
	}

	// create and link account TYPE to account
	switch accountType {
	case repository.AccountTypeSAVINGS:

		var rate pgtype.Numeric
		rate.Scan("0.4")
		savingsAccount, err := qtx.CreateSavingsAccount(a.ctx, repository.CreateSavingsAccountParams{
			AccountID:       account.ID,
			InterestRate:    rate,
			MinBalance:      100,
			WithdrawalLimit: 1000000,
		})
		if err != nil {
			return repository.Account{}, err
		}

		fmt.Printf("New SAVINGS ACCOUNT NUMBER: %s has been opened for USER %s\n", account.AccountNumber, user.FirstName+user.LastName)
		fmt.Printf("INFO: %+v\n", savingsAccount)

	case repository.AccountTypeCHECKING:
		checkingAccount, err := qtx.CreateCheckingAccount(a.ctx, repository.CreateCheckingAccountParams{
			AccountID:      account.ID,
			OverdraftLimit: 500,
			MaintenanceFee: 12,
		})
		if err != nil {
			return repository.Account{}, err
		}

		fmt.Printf("New CHECKING ACCOUNT NUMBER: %s has been openedd for USER %s\n", account.AccountNumber, user.FirstName+user.LastName)
		fmt.Printf("INFO: %+v\n", checkingAccount)
	}

	return account, tx.Commit(a.ctx)
}

func Shit() { fmt.Println("eat shit") }
