package discfac

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type DiscordFacadeImpl struct {
	guildID           string
	availableChannels []*discordgo.Channel
	mainChannel       *discordgo.Channel
	session           *discordgo.Session
}

func NewDiscordFacadeImpl(guildID string, session *discordgo.Session) *DiscordFacadeImpl {
	return &DiscordFacadeImpl{guildID: guildID, session: session}
}

func (fac *DiscordFacadeImpl) checkGuildChannels(channelName string) error {
	channels, err := fac.session.GuildChannels(fac.guildID)
	if err != nil {
		return err
	}

	fac.availableChannels = append(fac.availableChannels, channels...)
	for _, channel := range channels {
		if strings.EqualFold(channel.Name, channelName) {
			fac.mainChannel = channel
			return nil
		}
	}
	return nil
}

func (fac *DiscordFacadeImpl) InitMessageChannel(channelName string) error {
	if fac.mainChannel != nil {
		return nil
	}

	// Check firstly if the guild already has the requested channel
	if err := fac.checkGuildChannels(channelName); err != nil || fac.mainChannel != nil {
		return err
	}

	newChannel, err := fac.session.GuildChannelCreate(
		fac.guildID, channelName,
		discordgo.ChannelTypeGuildText,
	)
	if err != nil {
		return err
	}
	fac.mainChannel = newChannel

	return nil
}

func (fac *DiscordFacadeImpl) SendMessage(content string) error {
	_, err := fac.session.ChannelMessageSend(fac.mainChannel.ID, content)
	return err
}
