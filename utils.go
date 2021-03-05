package RadishCat

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

func (rc *RadishCat) GetUserMention(s string) (user *discordgo.User, err error){
	return rc.Session.User(strings.Replace(strings.Replace(s, "<@!", "", 1), ">", "", 1))
}