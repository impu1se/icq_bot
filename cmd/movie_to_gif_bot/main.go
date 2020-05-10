package main

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"

	"github.com/impu1se/icq_bot/configs"
	"github.com/impu1se/icq_bot/internal/botapi"
	"github.com/impu1se/icq_bot/internal/gif_bot"
	"github.com/impu1se/icq_bot/internal/storage"
	"go.uber.org/zap"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			if er, ok := r.(error); ok {
				fmt.Println(fmt.Sprintf("error: %v", er.Error()))
			}
			fmt.Println(fmt.Sprintf("stack_trace %v", debug.Stack()))
		}
	}()
	config := configs.NewConfig()

	botApi, err := botapi.NewBotApi(config)
	if err != nil {
		log.Fatalf("can't get new bot api, reason: %v", err)
	}

	db, err := storage.NewDb(config)
	if err != nil {
		log.Fatalf("can't create db, reason: %v", err)
	}

	ctx, _ := context.WithCancel(context.Background())
	logger := zap.NewExample()
	system := storage.NewLoader(logger)
	gifBot := gif_bot.NewGifBot(config, botApi.GetUpdatesChannel(ctx), system, db, logger, *botApi, ctx)

	gifBot.Run()

}
