package configuration

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync"
)

// Configuration holds configurable information
type Configuration struct {
	DbDriverName              string
	DbConnectionString        string
	JwtSecret                 string
	SMTP                      SMTP
	DentistProfileActivateURI string
	PatientProfileActivateURI string
	InvitationActivateURI     string
}

// SMTP holds smpt configurable information
type SMTP struct {
	Host     string
	Port     string
	Sender   string
	Password string
}

var instantiated *Configuration
var once sync.Once

// GetInstance implements singleton pattern
func GetInstance() *Configuration {
	once.Do(func() {
		file, _ := os.Open("conf.json")
		defer file.Close()
		decoder := json.NewDecoder(file)
		configuration := Configuration{}
		err := decoder.Decode(&configuration)
		if err != nil {
			fmt.Println("error:", err)
		}

		var dbDriverName string
		flag.StringVar(&dbDriverName, "dbDriverName", "", "specify 'dbDriverName' to use. Defaults to specified in config.json.")

		var dbConnectionString string
		flag.StringVar(&dbConnectionString, "dbConnectionString", "", "specify 'dbConnectionString' to use. Defaults to specified in config.json.")

		var jwtSecret string
		flag.StringVar(&jwtSecret, "jwtSecret", "", "specify 'jwtSecret' to use. Defaults to specified in config.json.")

		flag.Parse()

		if len(dbDriverName) > 0 {
			configuration.DbDriverName = dbDriverName
		}

		if len(dbConnectionString) > 0 {
			configuration.DbConnectionString = dbConnectionString
		}

		if len(jwtSecret) > 0 {
			configuration.JwtSecret = jwtSecret
		}

		instantiated = &configuration
	})

	return instantiated
}
