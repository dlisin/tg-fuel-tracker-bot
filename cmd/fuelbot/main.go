package main

import (
	"log"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot"
	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot/config"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("unable to load config", err)
	}

	app, err := bot.NewApp(cfg)
	if err != nil {
		log.Fatal("unable to create bot instance", err)
	}

	app.Run()
}
