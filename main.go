package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type config struct {
	Token         string   `json:"token"`
	AllowedGuilds []string `json:"allowed_guilds"`
	BotUserId     string   `json:"bot_user_id"`
}

func main() {
	context = appContext{}
	context.gameContext = &gameContext{}

	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Printf("Error reading configuration file!: %s\n", err.Error())
		return
	}

	context.config = &config{} // Parse configuration JSON
	err = json.Unmarshal(data, context.config)
	if err != nil {
		fmt.Printf("Invalid JSON parsed!: %s\n", err.Error())
		return
	}

	logIn(context.config)

	context.discord.AddHandler(messageCreate)

	reader := bufio.NewReader(os.Stdin)
	for running := true; running; {
		fmt.Printf("\nType 'help' or 'info', for more information.\nEnter Command: ")
		text := ReadStrippedString(reader)
		switch text {
		case "help", "info":
			fmt.Printf("\nType 'start', 'auto' to begin automatic battling\n" +
				"Type 'stop', 'exit' to exit the CLI and stop the Bot\n")
		case "start", "auto":
			fmt.Printf("\nStarting...\n")
			// TODO: Initiate battle sequence
		case "stop", "exit":
			fmt.Printf("\nStopping...\n")
			running = false
		default:
			fmt.Printf("\nInvalid Command: %s\n", text)
		}
	}

	// Cleanly close down the Discord session.
	err = context.discord.Close()
}

func logIn(botConfig *config) {
	// Create New Discord Session
	discord, err := discordgo.New(botConfig.Token)
	if err != nil {
		fmt.Printf("\nError occurred whilst creating Discord session!: %s\n", err.Error())
		os.Exit(1)
	}
	context.discord = discord

	err = context.discord.Open()
	if err != nil {
		fmt.Printf("\nError opening Discord Session!: %s\n", err.Error())
		err = context.discord.Close() // Close the session if it fails to open it
		os.Exit(1)
	}

}

func ReadStrippedString(reader *bufio.Reader) string {
	text, _ := reader.ReadString('\n')
	return strings.TrimSuffix(text, "\n")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID { // Remove this check for spam
		return
	}

	guildValid := false
	for _, guild := range context.config.AllowedGuilds {
		if m.GuildID == guild {
			guildValid = true
			break
		}
	}

	if !guildValid {
		fmt.Println("Not correct guild. Ignoring.")
		return
	}

	if m.Author.ID != context.config.BotUserId {
		fmt.Println("Message not sent by bot. Ignoring.")
		return
	}

	msg, err := s.ChannelMessages(m.ChannelID, 1, "", "", m.ID)
	if err != nil {
		fmt.Printf("UPDATE: Error getting message: %s\n", err.Error())
		return
	}

	if len(msg) > 0 {
		embeds := msg[0].Embeds

		// Reply to Anigame test command
		if msg[0].Content == "done." {
			_, _ = context.discord.ChannelMessageSend(m.ChannelID, strconv.Itoa(context.gameContext.currentLocation))
			_, _ = context.discord.ChannelMessageSend(m.ChannelID, strconv.Itoa(context.gameContext.currentFloor))
		}

		// Check for embeds
		if len(embeds) > 0 {
			fmt.Println(embeds[0].Title)

			if strings.Contains(embeds[0].Title, "Travelled to") { // Start Battle
				_, _ = context.discord.ChannelMessageSend(m.ChannelID, ".battle")
			}

			if embeds[0].Title == "**Victory <a:CHEER:705920932677681253>**" { // Progress to next Level if possible
				fmt.Println("Battle Won")
				_, _ = context.discord.ChannelMessageSend(m.ChannelID, ".fl next")
			}

			if strings.Contains(embeds[0].Description, "area ID you would like to go to.") { // Progress to next loc
				context.gameContext.currentLocation += 1
				_, err := context.discord.ChannelMessageSend(m.ChannelID, ".loc "+strconv.Itoa(context.gameContext.currentLocation))
				if err != nil {
					fmt.Println("Invalid location number!")
					return
				}
			}

			if strings.Contains(embeds[0].Title, "Successfully travelled to") { // Battle if travelled to next loc
				_, _ = context.discord.ChannelMessageSend(m.ChannelID, ".battle")
			}

			if strings.Contains(embeds[0].Description, "You do not have enough stamina to proceed!") {
				fmt.Println("\nInsufficient Stamina! Waiting for 30 Minutes before battling again.")
				time.Sleep(30 * time.Minute) // wait for n Minutes
				fmt.Println("\nAssuming Stamina is full. Battling Resumed.")
				_, _ = context.discord.ChannelMessageSend(m.ChannelID, ".battle")
			}

			// When the bot initiates a battle, strip the Title and get the current Loc and Floor numbers
			if strings.Contains(embeds[0].Title, "Challenging Floor") {
				s := stripFormattingCharacters(embeds[0].Title)
				sr := strings.Split(s, " ")       // Split the string into segments on basis of whitespaces
				sgmt := strings.Split(sr[2], "-") // Divide the Loc and Floor number by removing the '-'

				// Error checking for Location and Floor
				context.gameContext.currentLocation, err = strconv.Atoi(sgmt[0])
				if err != nil {
					fmt.Printf("Current Location invalid!%s\n", err.Error())
					return
				}
				context.gameContext.currentFloor, err = strconv.Atoi(sgmt[1])
				if err != nil {
					fmt.Printf("Current Floor invalid!%s\n", err.Error())
					return
				}
			}
		}
	}
}

// stripFormattingCharacters removes formatting characters, such as the ones used to create bold or italic text,
// from input and returns the sanitised string.
func stripFormattingCharacters(input string) string {
	input = strings.ReplaceAll(input, "__", "")
	input = strings.ReplaceAll(input, "**", "")
	input = strings.ReplaceAll(input, "~~", "")

	return input
}
