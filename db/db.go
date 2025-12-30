package db

import (
	"gorm.io/driver/postgres"
  	"gorm.io/gorm"
	"fmt"
	"errors"
)

const dsn string = "host=localhost user=postgres password=mango dbname=banking port=5432 sslmode=disable TimeZone=Asia/Kolkata"
var db *gorm.DB = nil

func startDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func init() {
	if db != nil {
		return
	}

	newDB, err := startDB()
	if err != nil {
		fmt.Println("Failed to start DB")
		fmt.Errorf("%+w\n", err)
		panic(err)
	}

	fmt.Printf("this from the db package gng: %+v\n", newDB)
	db = newDB
}

func DO_SOMETHING() User {
	fmt.Println("AWW HELL NAWW ;)")

	var users []User
	db.Find(&users)
	for _, user := range users {
		fmt.Printf("%+v", user)
	}

	return users[0]
}

func OpenSavingsAccount() error {
	tx := db.Begin()

	defer func() {
		KILLME()
		if r := recover(); r != nil {
			
			tx.Rollback()
		} else {
			fmt.Println("NO RECOVERY WHAT THE FUCK <3")
		}
	}()

	user := DO_SOMETHING()

	query := "INSERT INTO accounts (account_number, account_type, fk_user_id) VALUES (?, ?, ?) RETURNING *"

	var result Account
	if err := tx.Raw(query, "128659038628186094", SavingsAccount, user.ID).Scan(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			fmt.Println("DUP GLITCH")
		} else {
			fmt.Println("COOKED ERROR")
		}
		panic(err)
	}

	fmt.Printf("results of za fruits: %+v\n", result)

	return tx.Commit().Error
}

func KILLME() {
	var results []Account
	db.Find(&results)
	for _, account := range results {
		fmt.Printf("ACCOUNT INFO !!! : %+v\n", account)
	}
}

