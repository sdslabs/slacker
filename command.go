package slacker

import (
	"reflect"
	"strings"

	allot "github.com/sdslabs/allot/pkg"
)

// CommandDefinition structure contains definition of the bot command
type CommandDefinition struct {
	Description       string
	Example           string
	AuthorizationFunc func(botCtx BotContext, request Request) bool
	Handler           func(botCtx BotContext, request Request, response ResponseWriter)
}

// NewBotCommand creates a new bot command object
func NewBotCommand(usage string, definition *CommandDefinition, isParameterizedCommand bool, includeChannelIds []string) BotCommand {
	command := allot.New(usage)
	return &botCommand{
		usage:                  usage,
		definition:             definition,
		command:                command,
		isParameterizedCommand: isParameterizedCommand,
		includeChannelIds:      includeChannelIds,
	}
}

// BotCommand interface
type BotCommand interface {
	Usage() string
	Definition() *CommandDefinition
	IsParameterizedCommand() bool

	ContainsChannel(channelId string) bool
	MsgContains(text string) bool
	Match(req string) (allot.MatchInterface, error)
	Matches(text string) bool
	Tokenize() []*allot.Token
	Parameters() []allot.Parameter
	Execute(botCtx BotContext, request Request, response ResponseWriter)
}

// botCommand structure contains the bot's command, description and handler
type botCommand struct {
	usage                  string
	definition             *CommandDefinition
	command                *allot.Command
	isParameterizedCommand bool
	includeChannelIds      []string
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

func (c *botCommand) ContainsChannel(channelId string) bool {
	if reflect.DeepEqual(c.includeChannelIds, defaultIncludeChannelIds) {
		return true
	}
	for _, chId := range c.includeChannelIds {
		if chId == channelId {
			return true
		}
	}
	return false
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
