package util

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type config struct {
	ServerPort string `yaml:"ServerPort"`
	DBName     string `yaml:"DBName"`
	Salt       string `yaml:"Salt"`
}

func checkConfigExists(filename string) bool {
	fmt.Println("Checking if basic config files exist")
	_, fileError := os.Stat(filename)
	if fileError != nil {
		return false
	} else {
		return true
	}
}

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
func createSessionTable(db *sql.DB) (status bool) {
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

func checkDbExists(dbName string, salt []byte) (status string, error error) {
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
		tableStaus := createSessionTable(db)

	} else {
		fmt.Println("Session Table found")
	}
	fmt.Printf("Result: %s\n", result)
	return status, err
}

func GetConfig(filename string) (confContent config, confError error) {
	fmt.Println("Reading config File")
	configFile, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error in your config found: %q\n", err)
		panic(err)
	}
	var runtimeConf config
	err = yaml.Unmarshal(configFile, &runtimeConf)
	if err != nil {
		panic(err)
	}
	if runtimeConf.ServerPort == "" {
		confError = fmt.Errorf("ERROR: NO ServerPort configured, please check your configuration %s", filename)
		panic(confError)
	}
	return runtimeConf, confError
}

func writeSalt(filename string, salt []byte) {
	fmt.Println("Writing Salt to config")
	saltString := base64.StdEncoding.EncodeToString(salt)
	var conf config
	configFile, _ := os.Open(filename)
	decodedconf := yaml.NewDecoder(configFile)
	if err := decodedconf.Decode(&conf); err != nil {
		panic(err)
	}
	configFile.Close()
	conf.Salt = saltString
	updatedConf, marshErr := yaml.Marshal(&conf)
	if marshErr != nil {
		panic(marshErr)
	}
	if err := os.WriteFile(filename, updatedConf, 0644); err != nil {
		panic(err)
	}
}

func generateSalt() (salt []byte, saltError error) {
	fmt.Println("Generating Salt")
	salt = make([]byte, 128)
	return salt, saltError
}

func getSalt(readvalue string, filename string) (salt []byte, saltError error) {
	fmt.Println("Reading Salt from config")
	if readvalue == "" {
		fmt.Println("No Salt found or empty")
		salt, saltError = generateSalt()
		writeSalt(filename, salt)
	} else {
		fmt.Println("Salt found")
		//convert string to byteslice for use
		salt = []byte(readvalue)
	}
	return salt, saltError
}

func Initialization(filename string) (configState []string, ConfigError error) {
	config := checkConfigExists(filename)
	if config {
		fmt.Println("Config File found")
		serverConf, readconferr := GetConfig(filename)
		if readconferr != nil {
			fmt.Println(readconferr)
			panic(readconferr)
		}
		salt, saltErr := getSalt(serverConf.Salt, filename)
		if saltErr != nil {
			fmt.Println(saltErr)
			panic(saltErr)
		}
		DbFileName := fmt.Sprintf("%s.db", serverConf.DBName)
		dbState, createError := checkDbExists(DbFileName, salt)
		if createError != nil {
			panic(createError)
		}
		fmt.Printf("DB State: %s\n", dbState)

	} else {
		cwd, _ := os.Getwd()
		ConfigError := fmt.Errorf("Config File %q not found in %s", filename, cwd)
		panic(ConfigError)
	}
	return configState, ConfigError
}
