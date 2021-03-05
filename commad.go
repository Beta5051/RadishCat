package RadishCat

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"reflect"
	"strconv"
	"strings"
)

type Command struct {
	Name string
	Description string
	Aliases []string
	Handler CommandHandler
}

type CommandHandler interface {
	Run(rc *RadishCat, m *discordgo.Message) error
}

func NewCommand(name, description string, aliases []string, handler CommandHandler) *Command {
	t := reflect.TypeOf(handler)
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		return &Command{
			Name: name,
			Description: description,
			Aliases: aliases,
			Handler: handler,
		}
	}
	return nil
}

func (rc *RadishCat) HasCommand(name string) (int, bool){
	for idx, command := range rc.Commands {
		if command.Name == name {
			return idx, true
		}
	}
	return 0, false
}

func (rc *RadishCat) GetCommand(name string) *Command {
	if idx, ok := rc.HasCommand(name); ok{
		return rc.Commands[idx]
	}
	return nil
}

func (rc *RadishCat) AddCommand(command *Command) {
	if _, ok := rc.HasCommand(command.Name); !ok{
		rc.Commands = append(rc.Commands, command)
	}
}

func (rc *RadishCat) RemoveCommand(name string) {
	if idx, ok := rc.HasCommand(name); ok{
		rc.Commands = append(rc.Commands[:idx], rc.Commands[idx+1:]...)
	}
}

func (rc *RadishCat) commandMessageCreate(_ *discordgo.Session, m *discordgo.MessageCreate) {
	if rc.Session.State.User.ID == m.Author.ID {
		return
	}

	if !strings.HasPrefix(m.Content, rc.Prefix) {
		return
	}

	content := strings.Split(strings.TrimPrefix(m.Content, rc.Prefix), " ")
	for _, command := range rc.Commands {
		if command.Name == content[0] {
			err := rc.runCommand(command, m.Message, content[1:])
			if err != nil {
				rc.Logger.Errorln(err)
			}
			break
		} else {
			for _, alias := range command.Aliases {
				if alias == content[0] {
					err := rc.runCommand(command, m.Message, content[1:])
					if err != nil {
						rc.Logger.Errorln(err)
					}
					return
				}
			}
		}
	}
}

func (rc *RadishCat) RunCommand(name string, m *discordgo.Message, args []string) error {
	command := rc.GetCommand(name)
	if command != nil {
		return rc.runCommand(command, m, args)
	}
	return nil
}

func (rc *RadishCat) runCommand(command *Command, m *discordgo.Message, args []string) error {
	t := reflect.TypeOf(command.Handler).Elem()
	v := reflect.ValueOf(command.Handler).Elem()
	for i := 0; i < t.NumField() && i < len(args); i++ {
		f := v.Field(i)
		switch f.Kind() {
		case reflect.Int:
			iv, err := strconv.Atoi(args[i])
			if err != nil {
				return err
			}
			f.SetInt(int64(iv))
		case reflect.Uint:
			uv, err := strconv.ParseUint(args[i], 10, 64)
			if err != nil {
				return err
			}
			f.SetUint(uv)
		case reflect.String:
			f.SetString(args[i])
		case reflect.Ptr:
			if _, ok := f.Interface().(*discordgo.User); ok{
				user, err := rc.GetUserMention(args[0])
				if err != nil {
					return err
				}
				f.Set(reflect.ValueOf(user))
			}
		default:
			return errors.New("type that doesn't fit")
		}
	}
	return command.Handler.Run(rc, m)
}

type helpCommand struct {}

func (*helpCommand) Run(rc *RadishCat, m *discordgo.Message) error {
	embed := rc.NewEmbed("도움말")
	for _, command := range rc.Commands {
		var aliases string
		if len(command.Aliases) != 0 {
			var s string
			for idx, alias := range command.Aliases {
				s += alias
				if idx != len(command.Aliases) - 1 {
					s += ", "
				}
			}
			aliases = fmt.Sprintf("(%s)", s)
		}
		embed.AddField(command.Name + " " + aliases, command.Description, false)
	}

	_, err := rc.SendMessage(m.ChannelID, embed)
	return err
}