package bot

import (
	"MuscleBot/config" //importing config package
	"encoding/binary"
	"fmt"
	"io"
	"os"        //Used to monitor terminal for key enter
	"os/signal" //Used to determine what signal to give
	"syscall"   //used to call methods on the system such as kill
	"time"

	"github.com/bwmarrin/discordgo" //discordgo package from the repo of bwmarrin
	"github.com/jonas747/dca"       //Jonas747 dca encoding library
)

var BotId string               //Stores id of bot after making it a user
var GoBot *discordgo.Session   //goBot of type discordgo.Session
var buffer = make([][]byte, 0) //Used to store audio data

func Start() {
	GoBot, err := discordgo.New("Bot " + config.Token)

	//Error handling in case discord api fails
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := GoBot.User("@me") //Making bot a user

	//Error handling in case discord api fails
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//Storing our bot id given from api
	BotId = u.ID

	//Add Ready to the handler for call backs
	GoBot.AddHandler(Ready)

	//Adds handler function of our making (MessageHandler) from discord api
	GoBot.AddHandler(MessageHandler)

	//Adds handler function for creating guild for discordbot
	GoBot.AddHandler(GuildCreate)

	//Enables bot to identify which channel user is in that called for the bot
	GoBot.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	err = GoBot.Open()

	//Error handling in case discord api fails
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//Listens to console and waits for user to press contrl + c to terminate program and close the bot
	fmt.Println("Bot is up! Press control + c to terminate")
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-terminate

	//Closes down discord bot nicely
	GoBot.Close()
}

//Takes discordgo api's session and message that occurs within Discord
func MessageHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	//Apply check so bot doesn't reply to itself
	if message.Author.ID == BotId {
		return
	}

	//If user throws ping message at bot it'll respond with pong back in the Discord session
	//_ is a blank identifier that lets us get around needing to use variables within a func declaration
	if message.Content == "!ping" {
		_, _ = session.ChannelMessageSend(message.ChannelID, "pong")
	}

	if message.Content == "!MuscleBot" {
		session.ChannelMessageSend(message.ChannelID, "Joining")

		// Find the channel that the message came from.
		channel, err := session.State.Channel(message.ChannelID)
		if err != nil {
			// Could not find channel.
			return
		}

		// Find the guild for that channel.
		guild, err := session.State.Guild(channel.GuildID)
		if err != nil {
			// Could not find guild.
			return
		}

		// Look for the message sender in that guild's current voice states.
		for _, vs := range guild.VoiceStates {
			//If user ID matches the message maker's ID join that voice channel
			if vs.UserID == message.Author.ID {
				err = playSound(session, guild.ID, vs.ChannelID)
				if err != nil {
					fmt.Println(err.Error())
				}
				return
			}
		}
	}

}

//Function listens and responds when discord api sends ready event to bot
func Ready(session *discordgo.Session, event *discordgo.Ready) {
	//Sets playing status
	session.UpdateGameStatus(0, "!MuscleBot")
}

//Function called everytime a new guild (voice channel) is joined
func GuildCreate(session *discordgo.Session, event *discordgo.GuildCreate) {

	//If the voice channel can't be joined, exit
	if event.Guild.Unavailable {
		return
	}

	//Searches all channels and outputs this message to each channel
	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			//Plays on main text channel
			_, _ = session.ChannelMessageSend(channel.ID, "MuscleBot is ready! Type !MuscleBot to get a sound.")
			return
		}
	}
}

//Loads a sound from our sound files
func loadSound() {

	//Dca file options
	opts := dca.StdEncodeOptions
	//Stores File path 
	FileName := "MuscleBot/BotInfo/MuscleMan.mp3"

	// Encodes the file from mp3 to dca 
	encodeSession, err := dca.EncodeFile(FileName, opts)

	//Handles error incase encoding doesn't work
	if err != nil {
		return
	}

	// Make sure everything is cleaned up
	defer encodeSession.Cleanup()

	//Create MuscleMan.dca file from endcoded file
	output, err := os.Create("MuscleMan.dca")
	if err != nil {
		// Handle the error
		return
	}

	//Copy the encodeSession into output which stores the file data
	io.Copy(output, encodeSession)

	file, err := os.Open("MuscleMan.dca")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//Used to transfer bit data from dca file to binary read
	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return
			}
			return
		}

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}

// Plays buffer stored audio in the discord channel
func playSound(session *discordgo.Session, guildID, channelID string) (err error) {

	// Join the provided voice channel.
	vc, err := session.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Play Sound
	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	// Stop playing audio
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}
