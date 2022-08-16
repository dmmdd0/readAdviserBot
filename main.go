package main

import (
	"flag"
	"log"
	"readAdviserBot/clients/telegram"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	// token = flag.Get(token)
	//t := mustToken()

	// tgClient = telegram.New(token)
	// lesson 3
	tgClient := telegram.New(tgBotHost, mustToken())

	// receve new mesages from chat by API
	// lesson #4
	// fetcher = fetcher.New()

	// send mesages by API
	// lesson #4
	// processor = processor.New()

	// consumer.Srart(fetcher, processor)
}

// not GetToken, because Get haven't useful information
// must* for obligatory functions
// while parsing config only
func mustToken() string {
	// bot -tg-bot-token 'my token'
	// always write USAGE, for remember for
	token := flag.String(
		"token-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		// without token can't start anyway so os.Exit(1)
		log.Fatal("token is not specified")
	}

	return *token
}
