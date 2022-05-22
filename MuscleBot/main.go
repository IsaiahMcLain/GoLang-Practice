/*Discord Bot that is meant to join a voice channel and play Muscle Man sound bites from
Regular Show, bot was meant to learn go and working with an API, while it accomplished that
the bot is less than functional and will join voice channel but not play audio. The ping command
does work as intended. May come back to this project later on. Special thanks to the discordgo library 
and Jonas747 dca library */

package main

import (
	"MuscleBot/bot"
	"MuscleBot/config"
	"fmt"
)

func main() {
	err := config.ReadConfig() //Variable that reads config file and holds it's value

	//If there's an error reading the file, report it and end program
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Start()

	//Creating a channel that will hold a struct
	<-make(chan struct{})
	return

}
