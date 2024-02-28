package initialisation

import (
	"fmt"
	"os"
	conf "service_manager/init/config"
	crypt "service_manager/init/crypt"
	db "service_manager/init/db"
)

func Initialization(filename string) (configState []string, ConfigError error) {
	config := conf.CheckConfigExists(filename)
	if config {
		fmt.Println("Config File found")
		serverConf, readconferr := conf.GetConfig(filename)
		if readconferr != nil {
			fmt.Println(readconferr)
			panic(readconferr)
		}
		salt, saltErr := crypt.GetSalt(serverConf.Salt, filename)
		if saltErr != nil {
			fmt.Println(saltErr)
			panic(saltErr)
		}
		DbFileName := fmt.Sprintf("%s.db", serverConf.DBName)
		dbState, createError := db.CheckDbExists(DbFileName, salt)
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
