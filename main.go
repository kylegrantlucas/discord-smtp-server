package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/emersion/go-smtp"
	"github.com/kylegrantlucas/discord-smtp-server/email"
)

var discord *discordgo.Session

type backend struct{}

// Login handles a login command with username and password.
func (bkd *backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	if username != os.Getenv("SMTP_USERNAME") || password != os.Getenv("SMTP_PASSWORD") {
		return nil, errors.New("Invalid username or password")
	}
	return &session{}, nil
}

// AnonymousLogin requires clients to authenticate using SMTP AUTH before sending emails
func (bkd *backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return nil, smtp.ErrAuthRequired
}

// A Session is returned after successful login.
type session struct {
	webhook string
	from    string
}

func (s *session) Mail(from string, opts smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *session) Rcpt(to string) error {
	address, err := email.Parse(to)
	if err != nil {
		return err
	}

	guildID, err := getGuildID(address.TLD)
	if err != nil {
		return err
	}

	channelID, err := getChannelID(*guildID, address.Domain)
	if err != nil {
		return err
	}

	webhook, err := getWebhook(address.User, *channelID)
	if err != nil {
		return err
	}

	s.webhook = *webhook

	return nil
}

func (s *session) Data(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	reqBody, err := json.Marshal(map[string]string{
		"content": string(b),
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(
		s.webhook,
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (s *session) Reset() {}

func (s *session) Logout() error {
	return nil
}

func getGuildID(guildName string) (*string, error) {
	guilds, err := discord.UserGuilds(50, "", "")
	if err != nil {
		log.Print(err)
	}

	for _, guild := range guilds {
		if strings.ReplaceAll(strings.ToLower(guild.Name), " ", "") == strings.ToLower(guildName) {
			return &guild.ID, nil
		}
	}

	return nil, err
}

func getChannelID(guildID, channelName string) (*string, error) {
	channels, err := discord.GuildChannels(guildID)
	if err != nil {
		return nil, err
	}

	for _, channel := range channels {
		if strings.ToLower(channel.Name) == strings.ToLower(channelName) {
			return &channel.ID, nil
		}
	}

	log.Print("failed to find channel")

	return nil, err
}

func getWebhook(username, channelID string) (*string, error) {
	webhooks, err := discord.ChannelWebhooks(channelID)
	if err != nil {
		return nil, err
	}

	for _, hook := range webhooks {
		if strings.ToLower(hook.Name) == strings.ToLower(username) {
			webhook := fmt.Sprintf(
				"https://discord.com/api/webhooks/%v/%v",
				hook.ID,
				hook.Token,
			)
			return &webhook, nil
		}
	}

	log.Print("failed to build webhook")

	return nil, err
}

func main() {
	var err error
	be := &backend{}
	discord, err = discordgo.New(os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	s := smtp.NewServer(be)

	s.Addr = ":1025"
	s.Domain = "localhost"
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	log.Println("Starting server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
