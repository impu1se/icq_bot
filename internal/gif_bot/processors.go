package gif_bot

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/impu1se/icq_bot/internal/storage"

	"github.com/impu1se/icq_bot/configs"
	tgbotapi "github.com/mail-ru-im/bot-golang"
)

type System interface {
	Download(filepath, url string) error
	CreateNewDir(chatId string) error
	MakeGif(chatId string, dest string) error
	MakeImagesFromMovie(user *storage.User) error
	ClearDir(pattern string) error
}

type GifBot struct {
	Config  *configs.Config
	Updates <-chan tgbotapi.Event
	db      *storage.Database
	system  System
	logger  *zap.Logger
	ctx     context.Context
	api     tgbotapi.Bot
}

func NewGifBot(
	config *configs.Config,
	updates <-chan tgbotapi.Event,
	system System,
	db *storage.Database,
	logger *zap.Logger,
	api tgbotapi.Bot,
	ctx context.Context,
) *GifBot {
	return &GifBot{
		Config:  config,
		Updates: updates,
		system:  system,
		db:      db,
		logger:  logger,
		ctx:     ctx,
		api:     api,
	}
}

func (bot *GifBot) Run() {
	for update := range bot.Updates {

		if len(update.Payload.Parts) > 0 {
			bot.handlerVideo(&update)
			continue
		}

		if update.Payload.Message() != nil {
			bot.handlerMessages(&update)
			continue
		}
	}
}

func (bot *GifBot) NewMessage(chatId, message string, button *tgbotapi.ButtonResponse) error {

	if message == "" {
		return nil
	}
	text, err := bot.db.GetText(bot.ctx, message)
	if err != nil {
		bot.logger.Error(fmt.Sprintf("can't get text from db for message: %v", message))
		return err
	}
	msg := bot.api.NewTextMessage(chatId, text)
	if err := msg.Send(); err != nil {
		return err
	}
	return nil
}
