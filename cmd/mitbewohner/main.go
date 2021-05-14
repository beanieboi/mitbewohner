package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/beanieboi/mitbewohner/pkg/raidstatus"
	"github.com/bwmarrin/discordgo"
)

var (
	Token                  string
	NotificationsChannelID string
)

func main() {
	NotificationsChannelID, exists := os.LookupEnv("NOTIFICATION_CHANNEL_ID")

	if !exists {
		panic("please set NOTIFICATION_CHANNEL_ID")
	}

	Token, exists := os.LookupEnv("DISCORD_BOT_TOKEN")

	if !exists {
		panic("please set DISCORD_BOT_TOKEN")
	}

	bot, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	bot.Identify.Intents = discordgo.IntentsGuildMessages

	// Initialize our plugins
	bot.AddHandler(raidstatus.RaidStatusHandler)
	raidstatus.Initialize(bot, NotificationsChannelID)

	err = bot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	bot.Close()
}
