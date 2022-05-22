package config

import (
	"encoding/json"
	"fmt"       //Used to print errors
	"io/ioutil" //Helps read config.json file
)

var (
	Token     string //Stores auth token from config.json
	BotPrefix string //Stores bots prefix to be used

	config *ConfigStruct //Ptr to struct that stores values extracted from config.json
)

//Note: SPACING matters when it comes to the 'json:"Token" and other parts
type ConfigStruct struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
}

func ReadConfig() error {
	fmt.Println("Reading config file...")
	//ioutil function ReadFile reads json file and returns it's value into err, telling us if read was successful or not. It stores the read data into the file var
	file, err := ioutil.ReadFile("./config.json")

	//Handling error and printing it then returning it to main as well
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	//Print value of file var by explicitly converting to string
	fmt.Println(string(file))

	//Copying data from file into our Config Struct that was declared in var
	err = json.Unmarshal(file, &config)

	//Error handling for if Unmarshal fails
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	//Take data that was stored in ConfigStruct (config in var) and store the data into our declared variables
	Token = config.Token
	BotPrefix = config.BotPrefix

	//No error, return nil 
	return nil
}
