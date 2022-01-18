package slacker

import (
	"strings"

	allot "github.com/sdslabs/allot/pkg"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

// CommandDefinition structure contains definition of the bot command
type CommandDefinition struct {
	Description       string
	Examples          []string
	BlockID           string
	AuthorizationFunc func(botCtx BotContext, request Request) bool
	Handler           func(botCtx BotContext, request Request, response ResponseWriter)
	Interactive       func(*Slacker, *socketmode.Event, *slack.InteractionCallback, *socketmode.Request)

	// HideHelp will cause this command to not be shown when a user requests
	// help.
	HideHelp bool
}

// NewBotCommand creates a new bot command object
func NewBotCommand(usage string, definition *CommandDefinition, isParameterizedCommand bool) BotCommand {
	command := allot.New(usage)
	return &botCommand{
		usage:                  usage,
		definition:             definition,
		command:                command,
		isParameterizedCommand: isParameterizedCommand,
	}
}

// BotCommand interface
type BotCommand interface {
	Usage() string
	Definition() *CommandDefinition
	IsParameterizedCommand() bool

	MsgContains(text string) bool
	Match(req string) (allot.MatchInterface, error)
	Matches(text string) bool
	Tokenize() []*allot.Token
	Parameters() []allot.Parameter
	Execute(botCtx BotContext, request Request, response ResponseWriter)
	Interactive(*Slacker, *socketmode.Event, *slack.InteractionCallback, *socketmode.Request)
}

// botCommand structure contains the bot's command, description and handler
type botCommand struct {
	usage                  string
	definition             *CommandDefinition
	command                *allot.Command
	isParameterizedCommand bool
}

// Usage returns the command usage
func (c *botCommand) Usage() string {
	return c.usage
}

// Description returns the command description
func (c *botCommand) Definition() *CommandDefinition {
	return c.definition
}

// IsParameterizedCommand returns whether command is parameterized command or we only want substring match
func (c *botCommand) IsParameterizedCommand() bool {
	return c.isParameterizedCommand
}

func (c *botCommand) MsgContains(text string) bool {
	return strings.Contains(strings.ToLower(text), c.usage)
}

// Match determines whether the bot should respond based on the text received
func (c *botCommand) Match(text string) (allot.MatchInterface, error) {
	return c.command.Match(text)
}

// Matches checks if a comand definition matches a request
func (c *botCommand) Matches(text string) bool {
	return c.command.Matches(text)
}

// Tokenize returns the command format's tokens
func (c *botCommand) Tokenize() []*allot.Token {
	return c.command.Tokenize()
}

// Parameters returns the command format's tokens
func (c *botCommand) Parameters() []allot.Parameter {
	return c.command.Parameters()
}

// Execute executes the handler logic
func (c *botCommand) Execute(botCtx BotContext, request Request, response ResponseWriter) {
	if c.definition == nil || c.definition.Handler == nil {
		return
	}
	c.definition.Handler(botCtx, request, response)
}

// Interactive executes the interactive logic
func (c *botCommand) Interactive(slacker *Slacker, evt *socketmode.Event, callback *slack.InteractionCallback, req *socketmode.Request) {
	if c.definition == nil || c.definition.Interactive == nil {
		return
	}
	c.definition.Interactive(slacker, evt, callback, req)
}
