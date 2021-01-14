package discord

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Client struct {
	client *discordgo.Session
}

func NewClient(token string) (*Client, error) {
	client, err := discordgo.New(token)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) GetGuildID(guildName string) (*string, error) {
	guilds, err := c.client.UserGuilds(50, "", "")
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

func (c *Client) GetChannelID(guildID, channelName string) (*string, error) {
	channels, err := c.client.GuildChannels(guildID)
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

func (c *Client) GetWebhook(username, channelID string) (*string, error) {
	webhooks, err := c.client.ChannelWebhooks(channelID)
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
