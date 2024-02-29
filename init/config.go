package initialisation

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	ServerPort string `yaml:"ServerPort"`
	DBName     string `yaml:"DBName"`
	Salt       string `yaml:"Salt"`
}

func CheckConfigExists(filename string) Config {
	fmt.Println("Checking if basic config files exist")
	_, fileError := os.Stat(filename)
	if fileError != nil {
		panic(fileError)
	} else {
		serverConf, readconferr := GetConfig(filename)
		if readconferr != nil {
			fmt.Println(readconferr)
			panic(readconferr)
		}
		return serverConf
	}
}

func GetConfig(filename string) (confContent Config, confError error) {
	fmt.Println("Reading config File")
	configFile, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error in your config found: %q\n", err)
		panic(err)
	}
	var runtimeConf Config
	err = yaml.Unmarshal(configFile, &runtimeConf)
	if err != nil {
		panic(err)
	}
	if runtimeConf.ServerPort == "" {
		confError = fmt.Errorf("ERROR: NO ServerPort configured, please check your configuration '%s'", filename)
		panic(confError)
	}
	return runtimeConf, confError
}
