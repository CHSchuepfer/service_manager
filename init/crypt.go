package initialisation

import (
	"encoding/base64"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

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

func GetSalt(readvalue string, filename string) (salt []byte, saltError error) {
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
