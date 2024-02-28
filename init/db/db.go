package initialisation

import (
	"database/sql"
	"fmt"
	"os"
)

func createDb(dbName string) (status string, createErr error) {
	fmt.Println("Creating new Database")
	confDB, CreateDberr := sql.Open("sqlite3", dbName)
	status = confDB.Stats().WaitDuration.String()
	defer confDB.Close()

	if CreateDberr != nil {
		CreateDBError := fmt.Errorf("Error creating Database %s", dbName, "Exact Return Code: %q", CreateDberr)
		panic(CreateDBError)
	}

	fmt.Println("Database Created")
	return status, createErr
}

func checkUserData(db *sql.DB, salt []byte) (status bool) {
	fmt.Println("Checking for User Data")
	_, checkErr := db.Exec(`SELECT EXISTS (SELECT 1 FROM sessions WHERE username='admin')`)
	if checkErr != nil {
		panic(checkErr)
	} else {
		fmt.Println("User Data found")
		status = true
	}
	return status
}
func createSessionTable(db *sql.DB, salt []byte) (status bool) {
	fmt.Println("creating session table")
	tableCreateResult, dbCreateErr := db.Exec(
		`CREATE TABLE sessions (
        id INTEGER PRIMARY KEY,
        username TEXT NOT NULL,
        passwords TEXT NOT NULL
    )`,
	)
	if dbCreateErr != nil {
		fmt.Printf("Error on Table creation: %q", dbCreateErr)
		panic(dbCreateErr)
	} else {
		fmt.Println("Table Created")
	}
	checkuser := checkUserData(db, salt)
	if checkuser {
		status = true
	} else {
		status = false
	}
	return status
}

func CheckDbExists(dbName string, salt []byte) (status string, error error) {
	fmt.Printf("Checking if DB %s exists\n", dbName)
	_, dbfileError := os.Stat(dbName)
	cwd, _ := os.Getwd()
	if dbfileError != nil {
		fmt.Printf("DB not found under %s\n", cwd)
		_, dbErr := createDb(dbName)
		if dbErr != nil {
			return status, dbErr
		}

	}
	fmt.Println("Database found, checking for consistency")
	db, err := sql.Open("sqlite3", dbName)
	defer db.Close()
	if err != nil {
		panic(err)
	}

	result, err := db.Exec(`SELECT EXISTS (SELECT 1 FROM sqlite_master WHERE type='table' AND name='session')`)
	if err != nil {
		fmt.Println(err)
		tableStaus := createSessionTable(db, salt)

	} else {
		fmt.Println("Session Table found")
	}
	fmt.Printf("Result: %s\n", result)
	return status, err
}
