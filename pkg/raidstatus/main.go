package raidstatus

import (
	"fmt"
	"time"

	"github.com/beanieboi/raidstatus/raid"
	"github.com/bwmarrin/discordgo"
)

func Initialize(bot *discordgo.Session, channelID string) {
	go func() {
		for range time.Tick(time.Minute) {
			status, err := raid.Status()

			if err != nil {
				bot.ChannelMessageSend(channelID, err.Error())
			} else {
				for _, s := range status {
					if s.Status != "Online" {
						faultyDevices := "Faulty Devices: "
						for _, f := range s.FaultyDevices {
							faultyDevices = faultyDevices + fmt.Sprintf(" %s", f.BSDName)
						}
						bot.ChannelMessageSend(channelID, fmt.Sprintf("Status of '%s' is '%s' \n%s", s.Name, s.Status, faultyDevices))
					}
				}
			}
		}
	}()
}

func RaidStatusHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "raidstatus" {
		for _, status := range raid.StatusString() {
			s.ChannelMessageSend(m.ChannelID, status)
		}
	}
}
