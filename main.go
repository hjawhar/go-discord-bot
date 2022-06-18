package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/olekukonko/tablewriter"
)

var dg *discordgo.Session
var counter = 0

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DISCORD_TOKEN := os.Getenv("DISCORD_TOKEN")

	fmt.Println("Initiating bot...")
	// Create a new Discord session using the provided bot token.
	dg, err = discordgo.New("Bot " + DISCORD_TOKEN)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	if m.Content == "cars" {
		str := &strings.Builder{}
		data := [][]string{
			{"BMW", "E36 M3", "2JZ engine swapped", "1995"},
			{"BMW", "F22 M235i", "Single turbo", "2014"},
			{"Honda", "CRV", "-", "2022"},
		}

		table := tablewriter.NewWriter(str)
		table.SetHeader([]string{"Vendor", "Model", "Specs", "Year"})

		for _, v := range data {
			table.Append(v)
		}
		table.Render() // Send output
		out := "```" + str.String() + "```"
		s.ChannelMessageSend(m.ChannelID, out)
	}
}
