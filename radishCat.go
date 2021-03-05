package RadishCat

import (
	"fmt"
	"github.com/Beta5051/NeisGo"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"strings"
)

type RadishCat struct {
	Session *discordgo.Session
	Neis *NeisGo.Neis
	Logger *logrus.Logger
	Prefix string
	EmbedColor int
	Commands []*Command
}

func New(token string, v ...interface{}) (rc *RadishCat, err error) {
	if !strings.HasPrefix(token, "Bot ") {
		token = "Bot " + token
	}

	session, err := discordgo.New(token)
	if err != nil {
		return
	}
	rc = &RadishCat{
		Session: session,
		Neis: NeisGo.New(""),
		Logger: logrus.New(),
		Prefix: DEFAULT_PREFIX,
		EmbedColor: DEFAULT_EMBED_COLOR,
	}

	if len(v) >= 1 {
		if prefix, ok := v[0].(string); ok{
			rc.Prefix = prefix
		}
	}
	if len(v) >= 2 {
		if color, ok := v[1].(int); ok{
			rc.EmbedColor = color
		}
	}

	rc.Session.AddHandler(rc.ready)
	rc.Session.AddHandler(rc.commandMessageCreate)
	rc.AddCommand(NewCommand("도움말", "도움말", []string{"help"}, &helpCommand{}))
	rc.AddCommand(NewCommand("급식", "<지역> <학교> <학교종류> <날짜(필수X)> - 급식 정보를 보여줍니다.", []string{"diet"}, &dietCommand{}))
	return
}

func (rc *RadishCat) Open() error {
	return rc.Session.Open()
}

func (rc *RadishCat) Close() error {
	return rc.Session.Close()
}

func (rc *RadishCat) ready(_ *discordgo.Session, _ *discordgo.Ready) {
	err := rc.Session.UpdateGameStatus(0, fmt.Sprintf("%s or %s", rc.Prefix+"도움말", rc.Prefix+"help"))
	if err != nil {
		rc.Logger.Errorln(err)
	}
	rc.Logger.Println("Ready!")
}