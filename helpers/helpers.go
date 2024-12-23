package helpers

import (
	"crypto/sha512"
	"encoding/hex"
	"log"
	"strings"

	"github.com/greeneg/todoer/model"
)

func CheckIsNotLocked(u model.User) bool {
	return u.Status != "locked"
}

func CheckUserPass(username, password string) bool {
	user, err := model.GetUserByUserName(username)
	if err != nil {
		return false
	}

	status := CheckIsNotLocked(user)
	if !status {
		return false
	}

	// get the password hash from the user so we can compare it
	pwHash := user.PasswordHash

	// now calculate the sha512 of the password and see if it matches
	sha := sha512.Sum512([]byte(password))
	newPwHash := hex.EncodeToString(sha[:])

	return pwHash == newPwHash // returns boolean based on equality
}

func EmptyUserPass(username, password string) bool {
	return strings.Trim(username, " ") == "" || strings.Trim(password, " ") == ""
}

func FatalCheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
