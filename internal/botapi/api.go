package botapi

import (
	"fmt"
	"log"

	"github.com/impu1se/icq_bot/configs"
	tgbotapi "github.com/mail-ru-im/bot-golang"
)

func NewBotApi(config *configs.Config) (*tgbotapi.Bot, error) {

	fmt.Println("Running bot...")
	bot, err := tgbotapi.NewBot(config.ApiToken)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Printf("Authorized on account %s", "gif_movie_bot")
	return bot, nil
}
