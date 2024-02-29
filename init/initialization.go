package initialisation

import (
	"fmt"
	"os"
)

func Initialization(filename string) (configState []string, ConfigError error) {
	config := CheckConfigExists(filename)
	if config {
		fmt.Println("Config File found")
		serverConf, readconferr := GetConfig(filename)
		if readconferr != nil {
			fmt.Println(readconferr)
			panic(readconferr)
		}
		salt, saltErr := GetSalt(serverConf.Salt, filename)
		if saltErr != nil {
			fmt.Println(saltErr)
			panic(saltErr)
		}
		//Build the DB file name
		DbFileName := fmt.Sprintf("%s.db", serverConf.DBName)
		dbState, createError := CheckDbExists(DbFileName, salt)
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
