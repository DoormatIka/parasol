package main

import (
	"log"
	"os"

	// "github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/bot"
)

type Bot struct {
	Ctx *bot.Context
}
func (bot *Bot) Ping(*gateway.MessageCreateEvent) (string, error) {
	return "Pong!", nil;
}
func (bot *Bot) Help(*gateway.MessageCreateEvent) (string, error) {
	return bot.Ctx.Help(), nil;
}

// To run, do `BOT_TOKEN="TOKEN HERE" go run .`
// If in windows:
// set BOT_TOKEN=TOKEN_HERE
// go run .
func main() {
	token := os.Getenv("BOT_TOKEN");
	if token == "" {
		log.Fatalln("No BOT_TOKEN given. Do `BOT_TOKEN=\"TOKEN HERE\" go run .`");
	}

	commands := &Bot{};
	bot.Run(token, commands, func(ctx *bot.Context) error {
		ctx.HasPrefix = bot.NewPrefix("!!");
		return nil;
	});
}
