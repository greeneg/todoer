package main

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"time"
)

func convertSqliteTimestamp(t string) string {
	sqlTimestampFormat := "2006-01-02T15:04:05Z"
	timeFormat := "2006-01-02 15:04:05"
	createTime, _ := time.Parse(sqlTimestampFormat, t)
	return createTime.Format(timeFormat)
}

func createAccount(accountName string, accountFullName string, passwd string) (User, error) {
	t, err := DB.Begin()
	if err != nil {
		errPrintln("Could not start DB transaction!" + string(err.Error()))
		return User{}, err
	}

	q, err := t.Prepare("INSERT INTO Users (UserName, FullName, PasswordHash) VALUES (?, ?, ?)")
	if err != nil {
		errPrintln("Could not prepare the DB query!" + string(err.Error()))
		return User{}, err
	}

	// take password and hash it
	hash := sha512.Sum512([]byte(passwd))
	passwdHash := hex.EncodeToString(hash[:])

	// get the org Id

	_, err = q.Exec(accountName, accountFullName, passwdHash)
	if err != nil {
		errPrintln("Cannot create user '" + accountName + "': " + string(err.Error()))
		return User{}, err
	}

	t.Commit()

	user, err := getAccountByName(accountName)
	if err != nil {
		errPrintln("Could not retrieve user account: " + string(err.Error()))
		return User{}, err
	}

	return user, nil
}

func getAccountByName(accountName string) (User, error) {
	rec, err := DB.Prepare("SELECT Id,UserName,FullName,Status,CreationDate FROM Users WHERE UserName = ?")
	if err != nil {
		errPrintln("Could not prepare the DB query: " + string(err.Error()))
		return User{}, err
	}

	user := User{}
	err = rec.QueryRow(accountName).Scan(
		&user.Id,
		&user.UserName,
		&user.FullName,
		&user.Status,
		&user.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			errPrintln("No such role found in DB: " + string(err.Error()))
			return User{}, nil
		}
		errPrintln("Cannot retrieve role from DB: " + string(err.Error()))
		return User{}, err
	}

	user.CreationDate = convertSqliteTimestamp(user.CreationDate)

	return user, nil
}

func getAccountStatus(account string) (bool, error) {
	t, err := DB.Begin()
	if err != nil {
		errPrintln("Could not start DB transaction: " + string(err.Error()))
		return false, err
	}

	q, err := DB.Prepare("SELECT * FROM Users WHERE UserName IS ?")
	if err != nil {
		errPrintln("Could not prepare DB query! " + string(err.Error()))
		return false, err
	}

	err = q.QueryRow(account).Scan()
	if err != nil {
		if err != sql.ErrNoRows {
			errPrintln("Encountered error when querying database: " + string(err.Error()))
			return false, err
		}
		return false, err
	}

	t.Commit()

	return true, nil
}
