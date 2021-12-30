package pkg

import (
	"encoding/json"
	"log"
	"os"
)

//Configuration ...
type Configuration struct {
	Address      string
	ReadTimeout  int64
	WriteTimeout int64
	IdleTimeout  int64
	Static       string
}

//Config ...
var Config Configuration

//Logger ...
var Logger *log.Logger

//LoadConfig ...
func LoadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("Cannot open config file", err)
	}
	decoder := json.NewDecoder(file)
	Config = Configuration{}
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
}

//Info ...
func Info(args ...interface{}) {
	Logger.SetPrefix("INFO ")
	Logger.Println(args...)
}

//Danger ...
func Danger(args ...interface{}) {
	Logger.SetPrefix("ERROR ")
	Logger.Println(args...)
}

//Warning ..
func Warning(args ...interface{}) {
	Logger.SetPrefix("WARNING ")
	Logger.Println(args...)
}
