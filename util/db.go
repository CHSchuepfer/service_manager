package util

import (
	"database/sql"
	"fmt"
)

func genSecrets() {
	fmt.Println("in genSecrets")
}

func InitDB() {
	// Initialize the initRuntime
	sqliteDB, err := sql.Open("sqlite3", "./authager.db")
	if err != nil {
		fmt.Println(err)
	}
	defer sqliteDB.Close()
	// Create the table
	createTable := `
  CREATE TABLE IF NOT EXISTS users (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   username VARCHAR NOT NULL,
   password VARCHAR NOT NULL
   );
   `
	_, err = sqliteDB.Exec(createTable)
	if err != nil {
		fmt.Println(err)
	}

}
