package main

import (
	"context"
	"fmt"
	"log"
	"os"

	allot "github.com/sdslabs/allot/pkg"
	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"), slacker.WithDebug(true))
	bot.CustomCommand(func(usage string, definition *slacker.CommandDefinition) slacker.BotCommand {
		command := allot.New(fmt.Sprintf("custom-prefix %s", usage))
		return &cmd{
			usage:      usage,
			definition: definition,
			command:    command,
		}
	})

	// Invoked by `custom-prefix ping`
	bot.Command("ping", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			_ = response.Reply("it works!")
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

type cmd struct {
	usage      string
	definition *slacker.CommandDefinition
	command    *allot.Command
}

func (c *cmd) Usage() string {
	return c.usage
}

func (c *cmd) Definition() *slacker.CommandDefinition {
	return c.definition
}

func (c *cmd) Match(text string) (allot.MatchInterface, error) {
	return c.command.Match(text)
}

func (c *cmd) Matches(text string) bool {
	return c.command.Matches(text)
}

func (c *cmd) Tokenize() []*allot.Token {
	return c.command.Tokenize()
}

func (c *cmd) Parameters() []allot.Parameter {
	return c.command.Parameters()
}

func (c *cmd) Execute(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
	log.Printf("Executing command [%s] invoked by %s", c.usage, botCtx.Event().User)
	c.definition.Handler(botCtx, request, response)
}

func (c *cmd) Interactive(*slacker.Slacker, *socketmode.Event, *slack.InteractionCallback, *socketmode.Request) {
}
