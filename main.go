package main

import (
	"log"
	"os"
	"time"

	gosmtp "github.com/emersion/go-smtp"
	"github.com/kylegrantlucas/discord-smtp-server/smtp"
)

func main() {
	be, err := smtp.NewBackend(
		os.Getenv("DISCORD_TOKEN"),
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
	)
	if err != nil {
		log.Fatal(err)
	}

	s := gosmtp.NewServer(be)

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
