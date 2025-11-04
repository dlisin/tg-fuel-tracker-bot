package main

import (
	"log"

	"github.com/dlisin/tg-fuel-tracker-bot/internal/bot"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	app, err := bot.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	app.Run()
}
