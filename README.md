# Discord-SMTP Server
A simple relay that accepts SMTP messages and forwards them to a Discord webhook.

## Usage

### Local

```
env DISCORD_TOKEN=xxxxxxxxxxxx SMTP_USERNAME=username SMTP_PASSWORD=password go run main.go
```

### Docker

#### Run

```
docker run -t discord-smtp -e PORT=25 -e DISCORD_TOKEN=xxxxxxxxxxxx -e SMTP_USERNAME=username -e SMTP_PASSWORD=password
```

#### Compose

```
discord-smtp:
  image:
  container_name: discord-smtp
  env:
    - PORT=25
    - DISCORD_TOKEN=xxxxxxxxxxxx
    - SMTP_USERNAME=username
    - SMTP_PASSWORD=password
  restart: always
```

#### Testing

```
$ telnet localhost 1025
```

```
EHLO localhost
AUTH PLAIN
AHVzZXJuYW1lAHBhc3N3b3Jk
MAIL FROM:<test@test.com>
RCPT TO:<smtp@alert.karenplankton>
DATA
Hey
.
```

## Features

* SMTP Authentication
* Webhook Discovery