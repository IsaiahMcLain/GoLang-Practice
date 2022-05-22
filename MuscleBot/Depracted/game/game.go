// package game

// import (
// 	"math/rand"
// 	"time"

// 	"github.com/bwmarrin/discordgo" //discordgo package from the repo of bwmarrin
// )

// func GenerateString(Num int) string {
// 	//rune is used to convert between ascii and integer numbers, letters holds this rune
// 	var Letters = []rune("abcdefghijklmnopqrstuvwyxz")

// 	//Words will contain a random number of characters taken from Letters and it's size taken from what was passed in
// 	Words := make([]rune, Num)
// 	//Iterate over Words to create a random string
// 	for i := range Words {
// 		Words[i] = Letters[rand.Intn(len(Letters))]
// 	}

// 	return string(Words)
// }

// //Used to pass in size of string to GenerateString and returns the string given to it to the callee function
// func CreateRandString() string {
// 	//Seeds rand so that we get truly random val
// 	rand.Seed(time.Now().UnixNano())

// 	//Establish range size that word can be
// 	min := 5
// 	max := 25
// 	WordSize := rand.Intn(max-min+1) + min

// 	Word := GenerateString(WordSize)

// 	return Word

// }

// // //Makes discord bot wait 30 seconds to get user input to then check with if statement
// // func Wait() {
// // 	//s is seconds
// // 	for s := 0; s <= 30; s++ {
// // 		time.Sleep(time.Second * 1)
// // 	}
// // }

// func GamePlay(session *discordgo.Session, message *discordgo.MessageCreate) string {
// 	Input := message.Content
	
// 	return Input
// }

// func UserInput(Entry string, Word string) string {
// 	Win := "You Won!"
// 	TryAgain := "Try again!"
// 	if Entry == Word{
// 		return Win
// 	} 

// 	return TryAgain 
// }
