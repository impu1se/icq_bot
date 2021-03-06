package gif_bot

import (
	"fmt"
	"os"
	"runtime/debug"
	"strconv"

	"go.uber.org/zap"

	"github.com/impu1se/icq_bot/internal/storage"
	tgbotapi "github.com/mail-ru-im/bot-golang"
)

func (bot *GifBot) handlerMessages(update *tgbotapi.Event) {
	defer func() {
		if r := recover(); r != nil {
			if er, ok := r.(error); ok {
				fmt.Println(fmt.Sprintf("error: %v", er.Error()))
			}
			fmt.Println(fmt.Sprintf("stack_trace %v", debug.Stack()))
		}
	}()

	switch update.Payload.Text {
	case commandNewGif, clearTimes, oldGif:
		bot.handleNewGif(update)
	case start:
		bot.handleStart(update)
	case help:
		return
	default:
		bot.handleTimes(update)
	}
}

func (bot *GifBot) handlerVideo(update *tgbotapi.Event) {

	defer func() {
		if r := recover(); r != nil {
			if er, ok := r.(error); ok {
				fmt.Println(fmt.Sprintf("error: %v", er.Error()))
			}
			fmt.Println(fmt.Sprintf("stack_trace %v", debug.Stack()))
		}
	}()

	chatId := update.Payload.Chat.ID

	var fileId string
	for _, part := range update.Payload.Parts {
		if part.Type == "file" && part.Payload.Type == "video" {
			fileId = part.Payload.FileID
		}
	}
	if fileId == "" {
		if err := bot.NewMessage(chatId, "not video", nil); err != nil {
			bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
		}
		return
	}

	if err := bot.system.ClearDir(chatId); err != nil {
		bot.logger.Error("can't clear dir for new video")
		_ = bot.NewMessage(chatId, "error", nil)
		return
	}

	video, err := bot.api.GetFileInfo(fileId) // TODO: make check file size
	if err != nil {
		bot.logger.Error(fmt.Sprintf("can't get file from chat id: %v, reason: %v", chatId, err))
		if err := bot.NewMessage(chatId, "download error", nil); err != nil {
			bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
		}
		return
	} else {
		if err := bot.NewMessage(chatId, "save video", nil); err != nil {
			bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
		}
	}

	err = bot.system.Download(fmt.Sprintf("%v/%v", chatId, video.Name), video.URL)
	if err != nil {
		bot.logger.Error(fmt.Sprintf("can't download video, reason %v", err))
		_ = bot.NewMessage(chatId, "error", nil)
		return
	}
	if err := bot.db.UpdateLastVideo(bot.ctx, chatId, video.Name); err != nil {
		bot.logger.Error(fmt.Sprintf("can't update last video, reason %v", err))
		_ = bot.NewMessage(chatId, "error", nil)
		return
	}
	if err := bot.db.ClearTime(bot.ctx, chatId); err != nil {
		bot.logger.Error(fmt.Sprintf("can't clear time, reason %v", err))
		_ = bot.NewMessage(chatId, "error", nil)
		return
	}
	if err := bot.NewMessage(chatId, "successful download", nil); err != nil {
		bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
		return
	}
}

func (bot *GifBot) handleStart(update *tgbotapi.Event) {
	chatInfo, err := bot.api.GetChatInfo(update.Payload.Chat.ID)
	if err != nil {
		bot.logger.Error("can't get chat info :", zap.Field{String: err.Error()})
		return
	}
	user := &storage.User{
		ChatId:   update.Payload.Chat.ID,
		UserName: chatInfo.FirstName,
	}

	if err := bot.db.CreateUser(bot.ctx, user); err != nil {
		bot.logger.Error("can't crete user, reason:", zap.Field{String: err.Error()})
		_ = bot.NewMessage(user.ChatId, "error", nil)
		return
	}

	if err := bot.system.CreateNewDir(user.ChatId); err != nil {
		bot.logger.Error(fmt.Sprintf("can't create new dir for user with chat %v, reason %v", user.UserName, err))
	}

	if err := bot.NewMessage(user.ChatId, update.Payload.Text, []interface{}{user.UserName}); err != nil {
		bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
		return
	}
}

func (bot *GifBot) handleNewGif(update *tgbotapi.Event) {
	chatId := update.Payload.Chat.ID

	if err := bot.db.ClearTime(bot.ctx, chatId); err != nil {
		bot.logger.Error(fmt.Sprintf("can't clear time for user with id %v, reason: %v", chatId, err))
		_ = bot.NewMessage(chatId, "error", nil)
		return
	}

	if err := bot.NewMessage(chatId, update.Payload.Text, nil); err != nil {
		bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
		return
	}
}

func (bot *GifBot) handleTimes(update *tgbotapi.Event) {
	chatId := update.Payload.Chat.ID
	time, err := strconv.Atoi(update.Payload.Text)
	if err != nil {
		bot.logger.Error("can't parse time from message")
		if err := bot.NewMessage(chatId, "invalid message", nil); err != nil {
			bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
		}
		return
	}

	user, err := bot.db.GetUser(bot.ctx, chatId)
	if err != nil {
		bot.logger.Error(fmt.Sprintf("can't get user by chat id: %v, reason: %v", chatId, err))
		_ = bot.NewMessage(chatId, "error", nil)
		return
	}

	if user.StartTime == nil {
		if err := bot.db.UpdateStartTime(bot.ctx, chatId, time); err != nil {
			bot.logger.Error(fmt.Sprintf("can't update start time by chat id: %v, reason: %v", chatId, err))
			_ = bot.NewMessage(chatId, "error", nil)
			return
		}
		if err := bot.NewMessage(chatId, "end second", nil); err != nil {
			bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
			return
		}
	} else {
		if message, valid := checkValidTimes(time, *user.StartTime); !valid {
			if err := bot.NewMessage(chatId, message, nil); err != nil {
				bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
			}
			return
		}

		if err := bot.db.UpdateEndTime(bot.ctx, chatId, time); err != nil {
			bot.logger.Error(fmt.Sprintf("can't update end time by chat id: %v, reason: %v", chatId, err))
			_ = bot.NewMessage(chatId, "error", nil)
			return
		}
		endTime := time - *user.StartTime
		user.EndTime = &endTime
		if err := bot.NewMessage(chatId, "create video", nil); err != nil {
			return
		}

		err = bot.system.MakeImagesFromMovie(user)
		if err != nil {
			bot.logger.Error(fmt.Sprintf("can't make image from movie, reason: %v", err))
			return
		}
		if err := bot.NewMessage(chatId, "start create video", nil); err != nil {
			bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
			return
		}

		scale, err := bot.db.GetScale(bot.ctx)
		if err != nil {
			bot.logger.Error(fmt.Sprintf("can't get scale from db, reason: %v", err))
			scale = 0.5
		}
		gifPath := fmt.Sprintf("%v/%v.gif", chatId, user.LastVideo)
		err = bot.system.MakeGif(chatId, gifPath, scale)
		if err != nil {
			bot.logger.Error(fmt.Sprintf("can't make gif from movie, reason: %v", err))
			_ = bot.NewMessage(chatId, "error", nil)
			return
		}
		if err := bot.NewMessage(chatId, "loading gif", nil); err != nil {
			bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
			_ = bot.NewMessage(chatId, "error", nil)
			return
		}
		gif, err := os.Open("user_data/" + gifPath)
		if err != nil {
			bot.logger.Error(fmt.Sprintf("can't open file"))
			_ = bot.NewMessage(chatId, "error", nil)
			return
		}
		if err := bot.api.NewFileMessage(chatId, gif).Send(); err != nil {
			bot.logger.Error(fmt.Sprintf("can't send gif, reason: %v", err))
			_ = bot.NewMessage(chatId, "error", nil)
			return
		}
	}
}

func checkValidTimes(endTime, startTime int) (string, bool) {
	if endTime <= startTime {
		return "end more start", false
	} else if endTime-startTime > 10 {
		return "video more 10s", false
	}
	return "", true
}
