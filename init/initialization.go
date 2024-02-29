package initialisation

import (
	"fmt"
)

func Initialization(filename string) (configState []string, ConfigError error, salt []byte) {
	serverConfig := CheckConfigExists(filename)
	fmt.Println("Config File found")
	salt, saltErr := GetSalt(serverConfig.Salt, filename)
	if saltErr != nil {
		fmt.Println(saltErr)
		panic(saltErr)
	}
	//Build the DB file name
	DbFileName := fmt.Sprintf("%s.db", serverConfig.DBName)
	dbState, createError := CheckDbExists(DbFileName, salt)
	if createError != nil {
		panic(createError)
	}
	fmt.Printf("DB State: %s\n", dbState)
	return configState, ConfigError, salt
}
