package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/pborman/getopt/v2"
	"golang.org/x/term"
)

var VERSION string = "0.0.1"

var appFullPath, _ = os.Executable()
var app = path.Base(appFullPath)

// setup the global flags
var (
	dbFile     string
	account    string
	fullName   string
	optHelp    = getopt.BoolLong("help", 'h', "This help message")
	optVersion = getopt.BoolLong("version", 'v', "Show the version")
)

func showHelp() {
	println(app + " - Setup tool for Allocator Daemon")
	dividerLine := strings.Repeat("=", 43)
	println(dividerLine)
	println("Add and configure roles or accounts for the Allocator Daemon\n")
	println("OPTIONS:")
	println("   -d|--database-file FILENAME_PATH       REQUIRED: The full or relative path")
	println("                                          to the database file")
	println("   -a|--account ACCOUNT_NAME              OPTIONAL: The account to create")
	println("   -r|--role ROLE_NAME                    OPTIONAL: The role to create")
	println("   -o|--org-unit ORGANIZATIONAL_UNIT_NAME OPTIONAL: The organizational unit to")
	println("                                          create")
	println("   -f|--fullname QUOTED_FULLNAME          CONDITIONALLY OPTIONAL: If the")
	println("                                          account flag is set, this is required.")
	println("                                          This should be the full name or")
	println("                                          description for the account to be")
	println("                                          registered with the system.")
	println("   -D|--role-description ROLE_DESCRIPTION CONDITIONALLY OPTIONAL: If the")
	println("                                          role flag is set, this is required.")
	println("                                          This should be the description for")
	println("                                          the role to be registered with the")
	println("                                          system.")
	println("   -O|--org-unit-description DESCRIPTION  CONDITIONALLY OPTIONAL: If the")
	println("                                          org unit flag is set, this is required.")
	println("                                          This describes the organizational unit")
	println("                                          for the org-unit to be registered with")
	println("                                          the system.")
	println("")
	println("Author: Gary L. Greene, Jr. <greeneg@tolharadys.net>")
	println("License: Apache Public License, v2")
	showVersion()
}

func showVersion() {
	println("version: " + VERSION)
}

func init() {
	getopt.FlagLong(&dbFile, "database-file", 'd', "The full path to the database file")
	getopt.FlagLong(&account, "account", 'a', "The account to add to the system")
	getopt.FlagLong(&fullName, "fullname", 'f', "The full name to associate with the account")
}

func main() {
	getopt.Parse()

	if *optHelp {
		showHelp()
		os.Exit(0)
	}


	if *optVersion {
		showVersion()
		os.Exit(0)
	}


	if dbFile != "" {
		println("Database file: " + dbFile)
		ConnectDatabase(dbFile)
		infoPrintln("Database connection completed")
	} else {
		errPrintln("Database file must be defined.")
		showHelp()
		os.Exit(1)
	}

	// do we need to process an account the user passed in?
	if account != "" {
		println("Account: " + account)
		if fullName != "" {
			println("Account fullname: " + fullName)
		} else {
			errPrintln("Account must have a full name")
			showHelp()
			os.Exit(1)
		}
	}

	// check if account already exists
	accountStatus, err := getAccountStatus(account)
	if err != nil && err != sql.ErrNoRows {
		errPrintln("Encountered error when checking account status: " + string(err.Error()))
		os.Exit(1)
	}
	if !accountStatus {
		fmt.Print("Enter new password: ")
		input, _ := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Print("\nRe-enter passphrase: ")
		input2, _ := term.ReadPassword(int(os.Stdin.Fd()))
		println("")
		if strings.Compare(string(input), string(input2)) != 0 {
			errPrintln("Password does not match. Exiting")
			os.Exit(1)
		}
		accountRecord, err := createAccount(account, fullName, string(input))
		if err != nil {
			errPrintln("Encountered error when creating account '" + account + "': " + string(err.Error()))
			os.Exit(1)
		}
		accountRecordStr, err := json.Marshal(accountRecord)
		if err != nil {
			errPrintln("Encountered error when converting struct to JSON: " + string(err.Error()))
			os.Exit(1)
		}
		infoPrintln("account '" + account + "' created: " + string(accountRecordStr))

	}
}
