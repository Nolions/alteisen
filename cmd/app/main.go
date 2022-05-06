package main

import (
	"context"
	"errors"
	"flag"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nolions/alteisen/conf"
	"github.com/nolions/alteisen/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	confPath string // config path
)

func main() {
	flag.StringVar(&confPath, "c", "config.yaml", "default config path")
	flag.Parse()

	config, err := conf.New(confPath)
	if err != nil {
		log.Fatal(errors.New(err.Error()))
	}

	bot, err := tgbotapi.NewBotAPI(config.Bot.Token)
	if err != nil {
		log.Fatalf("Create telegram bot fail, errpr: %v\n", err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	log.Printf("Authorized's token %s", bot.Token)

	//config.Bot.BaseUrl

	webhookUrl := config.Bot.BaseUrl +"/webhook/" + bot.Token
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(webhookUrl))
	println(webhookUrl)
	if err != nil {
		log.Printf("set Webhook fail, error:%v\n", err)
	}

	ctx := context.Background()
	app := server.New(ctx, bot)
	serv := server.NewHttpServer(app, &config.HttpServ)
	serv.Run()

	shutdown(&config.App, serv)
}

func shutdown(conf *conf.App, srv *server.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	s := <-c
	log.Printf("get a signal %s. (%s) Server is shutting down ...\n", s.String(), conf.Name)

	// close http server with timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("http server shutdown: %v\n", err)
	}

	log.Printf("(%s) Server is exit.\n", conf.Name)
}
