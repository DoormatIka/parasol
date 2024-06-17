package main

import (
	"log"
	"os"
	"os/signal"
	"reflect"
	"regexp"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
);
var command_regex, err = regexp.Compile(`^!!(\w+)\s([\s\S+]+)$`);

type Commands struct {
	Ping Command
};
type Command struct {
	description string
	execute func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)
};

func (c *Commands) findCommand(name string) *Command {
	t := reflect.TypeOf(c).Elem();
	v := reflect.ValueOf(c).Elem();
	for i := 0; i < t.NumField(); i++ {
		method := t.Field(i);
		if strings.ToLower(method.Name) == name {
			cmd := v.Field(i).Interface().(Command);
			return &cmd;
		}
	}
	return nil;
}

func (c *Commands) RunCommand(sesh *discordgo.Session, msg *discordgo.MessageCreate) {
	println(msg.Content);
	matches := command_regex.FindStringSubmatch(msg.Content);
	if matches != nil {
		name := matches[1];
		args := strings.Split(matches[2], " ");
		cmd := c.findCommand(name);
		if cmd != nil {
			cmd.execute(sesh, msg, args);
		}
	}
}

// To run, do `BOT_TOKEN="TOKEN HERE" go run .`
// If in windows:
// set BOT_TOKEN=TOKEN_HERE
// go run .
func main() {
	if err != nil {
		log.Fatalln("Regex:", err);
	}
	commands := &Commands{
		Ping: Command{ 
			description: "Responds with pong.",
			execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
				_, err := s.ChannelMessageSendReply(m.ChannelID, "Pong.", m.Reference());
				if err != nil {
					println("Error: ", err);
				}
			},
		},
	};


	token := os.Getenv("BOT_TOKEN");
	if token == "" {
		log.Fatalln("No token provided.");
	}
	dg, err := discordgo.New("Bot " + token);
	if err != nil {
		log.Fatalln("Error opening connection,", err);
	}
	dg.AddHandler(func (s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.Bot {
			return;
		}
		commands.RunCommand(s, m);
	});
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

