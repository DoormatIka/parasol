package main

import (
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
);

var command_regex, err = regexp.Compile(`\?\?(\w+)\s(.+)`);

type Commands struct {};

func (p *Commands) RunCommand(s string) {
	matches := command_regex.FindStringSubmatch(`??quit args1 arg5`);
	if matches != nil {
		command := matches[1];
		args := strings.Split(matches[2], " ");
		// how do i get the command name and activate the ping command?
		// my idea rn is to make a map[string]func and just do a get(), but there has to be a better way?
	}
}

func (c *Commands) Ping(s *discordgo.Session, msg *discordgo.MessageCreate) {
	
}

// To run, do `BOT_TOKEN="TOKEN HERE" go run .`
// If in windows:
// set BOT_TOKEN=TOKEN_HERE
// go run .
func main() {
	if err != nil {
		log.Fatalln("Regex:", err);
	}

	token := os.Getenv("BOT_TOKEN");
	if token == "" {
		log.Fatalln("No token provided.");
	}
	dg, err := discordgo.New("Bot " + token);
	if err != nil {
		log.Fatalln("Error opening connection,", err);
	}
	dg.AddHandler(MessageCreate);
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentMessageContent;

	if err := dg.Open(); err != nil {
		log.Fatalln("Error oppening connection,", err);
	}
	log.Println("Bot is now running. Press CTRL-C to exit.");
	sc := make(chan os.Signal, 1);
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill);
	<-sc // blocks main loop until sc recieves something, which are these syscalls.

	// cleanly closes out the discord sesh
	dg.Close();
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return;
	}
	if m.Content == "!!ping" {
		_, err := s.ChannelMessageSendReply(m.ChannelID, "Pyon.", m.Reference());
		if err != nil {
			log.Println("Ping failed to send reply.", err);
		}
	}
}
