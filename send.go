package RadishCat

import (
	"errors"
	"github.com/bwmarrin/discordgo"
)

func (rc *RadishCat) SendMessage(channelID string, v interface{}) (m *discordgo.Message, err error){
	if s, ok := v.(string); ok{
		return rc.Session.ChannelMessageSend(channelID, s)
	} else if e, ok := v.(*Embed); ok{
		return rc.Session.ChannelMessageSendEmbed(channelID, e.MessageEmbed)
	} else {
		return nil, errors.New("type that doesn't fit")
	}
}

func (rc *RadishCat) SendMessageUser(user *discordgo.User, v ...interface{}) (m *discordgo.Message, err error) {
	channel, err := rc.Session.UserChannelCreate(user.ID)
	if err != nil {
		return
	}
	return rc.SendMessage(channel.ID, v)
}