package RadishCat

import (
	"github.com/bwmarrin/discordgo"
)

type Embed struct {
	*discordgo.MessageEmbed
}

func (rc *RadishCat) NewEmbed(title string) *Embed {
	return &Embed{
		MessageEmbed: &discordgo.MessageEmbed{
			Title: title,
			Color: rc.EmbedColor,
			Footer: &discordgo.MessageEmbedFooter{
				Text: rc.Session.State.User.Username,
				IconURL: rc.Session.State.User.AvatarURL(""),
			},
		},
	}
}

func (e *Embed) SetDescription(s string) *Embed {
	e.Description = s
	return e
}

func (e *Embed) AddField(name, value string, inline bool) *Embed {
	if e.Fields == nil {
		e.Fields = []*discordgo.MessageEmbedField{}
	}

	e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
		Name: name,
		Value: value,
		Inline: inline,
	})
	return e
}