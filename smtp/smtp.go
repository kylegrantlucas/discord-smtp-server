package smtp

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/emersion/go-smtp"
	"github.com/kylegrantlucas/discord-smtp-server/discord"
	"github.com/kylegrantlucas/discord-smtp-server/email"
)

type Backend struct {
	discordClient *discord.Client
	username      string
	password      string
}

func NewBackend(discordToken, username, password string) (*Backend, error) {
	discordClient, err := discord.NewClient(discordToken)
	if err != nil {
		return nil, err
	}

	return &Backend{
		discordClient: discordClient,
		username:      username,
		password:      password,
	}, nil
}

func (b *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	if username != b.username || password != b.password {
		return nil, errors.New("Invalid username or password")
	}
	return &Session{
		backend: b,
	}, nil
}

func (b *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return nil, smtp.ErrAuthRequired
}

type Session struct {
	backend *Backend
	webhook string
	from    string
}

func (s *Session) Mail(from string, opts smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string) error {
	address, err := email.Parse(to)
	if err != nil {
		return err
	}

	guildID, err := s.backend.discordClient.GetGuildID(address.TLD)
	if err != nil {
		return err
	}

	channelID, err := s.backend.discordClient.GetChannelID(*guildID, address.Domain)
	if err != nil {
		return err
	}

	webhook, err := s.backend.discordClient.GetWebhook(address.User, *channelID)
	if err != nil {
		return err
	}

	s.webhook = *webhook

	return nil
}

func (s *Session) Data(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	reqBody, err := json.Marshal(
		map[string]string{
			"content": string(b),
		},
	)
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

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}
