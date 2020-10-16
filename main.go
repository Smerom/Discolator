package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Smerom/Disclator/characterTracker"

	translate "cloud.google.com/go/translate/apiv3"
	"github.com/bwmarrin/discordgo"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
)

// Variables used for command line parameters
var (
	Token   string
	Dev     string
	Project string
	MaxChar int
)

var ct characterTracker.Tracker

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&Dev, "d", "", "Dev ID")
	flag.StringVar(&Project, "p", "", "Project ID")
	flag.IntVar(&MaxChar, "maxchar", 500000, "Max character count")
	flag.Parse()

	ct = characterTracker.NewMemoryTracker()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Printf("Message recieved: %s", m.Content)

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		log.Printf("From bot user")
		return
	}

	if m.Author.ID != Dev {
		log.Printf("Not from dev")
		return
	}

	// check translation count
	if ct.CountAfterString(m.Content) > MaxChar {
		log.Printf("Exceeds translation count")

		// message with count exceeded
		s.ChannelMessageSend(m.ChannelID, "Translation character count exceeded.")
		return
	}

	ctx := context.Background()
	c, err := translate.NewTranslationClient(ctx)
	if err != nil {
		// TODO: Handle error.
		s.ChannelMessageSend(m.ChannelID, "Client failed.")
		return
	}

	req := &translatepb.TranslateTextRequest{
		Contents:           []string{m.Content},
		SourceLanguageCode: "en",
		TargetLanguageCode: "ru",
		Parent:             fmt.Sprintf("projects/%s", Project),
	}

	// add character count to tracker
	ct.AddCharacters(m.Content)

	resp, err := c.TranslateText(ctx, req)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Translation failed: %s", err))
		return
	}

	log.Printf("Translated text: %s", resp.Translations[0].TranslatedText)

	s.ChannelMessageSend(m.ChannelID, resp.Translations[0].TranslatedText)
}
