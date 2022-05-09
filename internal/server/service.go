package server

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nolions/alteisen/conf"
	"log"
	"net/http"
	"time"
)

type Application struct {
	Ctx  context.Context
	Bot  *tgbotapi.BotAPI
	Conf *conf.Bot
}

type Server struct {
	HttpServer *http.Server
}

func New(ctx context.Context, bot *tgbotapi.BotAPI, config *conf.Bot) *Application {
	return &Application{
		Ctx: ctx,
		Bot: bot,
		Conf: config,
	}
}

// NewHttpServer
// init http server
func NewHttpServer(app *Application, conf *conf.HttpServ) *Server {
	e := engine()
	app.router(e)

	addr := fmt.Sprintf(":%s", conf.Addr)
	h := &http.Server{
		Addr:         addr,
		Handler:      e,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return &Server{
		HttpServer: h,
	}
}

// Run
// run http server
func (app *Server) Run() {
	if err := app.HttpServer.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

// Shutdown
// shut down service
func (srv *Server) Shutdown(ctx context.Context) error {
	return srv.Shutdown(ctx)
}
