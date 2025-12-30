package db

import (
	"database/sql/driver"
	"errors"
)

type User struct {
    ID          string
    FirstName   string
    LastName    string 
    Email       string 
    PhoneNumber string 
    Address     string 
    Password    string 
    CreatedAt   string 
    UpdatedAt   string 
}

/*
type User struct {
    ID          string `gorm:"column:id;primaryKey"`
    FirstName   string `gorm:"column:first_name"`
    LastName    string `gorm:"column:last_name"`
    Email       string `gorm:"column:email"`
    PhoneNumber string `gorm:"column:phone_number;not null;unique"`
    Address     string `gorm:"column:address"`
    Password    string `gorm:"column:password;not null"`
    CreatedAt   string `gorm:"column:created_at"`
    UpdatedAt   string `gorm:"column:updated_at"`
}
*/

type AccountTypeEnum string

const (
	SavingsAccount AccountTypeEnum = "SAVINGS"
	CheckingAccount AccountTypeEnum = "CHECKING"
)

var (
	AccountTypeToEnum = map[string]AccountTypeEnum{
		"SAVINGS": SavingsAccount,
		"CHECKING": CheckingAccount,
	}

	AccountEnumToType = map[AccountTypeEnum]string{
		SavingsAccount: "SAVINGS",
		CheckingAccount: "CHECKING",
	}
)

func (at *AccountTypeEnum) Scan(value interface{}) error {
	accountType, ok := value.(string)

	if !ok {
		return errors.New("idk wtf happened you on yo own with this one my g")
	}

	*at = AccountTypeToEnum[accountType]

	return nil
}

func (at *AccountTypeEnum) Value() (driver.Value, error) {
	accountString, ok := AccountEnumToType[*at]
	if !ok {
		return nil, errors.New("nah gang we cooked")
	}
	return accountString, nil
}

type Account struct {
	ID string
	AccountNumber string
	Balance int64
	AccountType AccountTypeEnum
	FkUserId string
	CreatedAt string
}

type Savings_Account struct {}
type Checking_Account struct {}
type Transactions struct {}


