package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
)

func (app Application) router(e *gin.Engine) {
	e.HandleMethodNotAllowed = true
	e.NoMethod(HandleNoAllowMethod)
	e.NoRoute(HandleNotFound)

	e.GET("/health", ErrHandler(app.healthHandler))
	e.POST("/webhook/"+app.Bot.Token, ErrHandler(app.webhookHandler))
}

func (app Application) healthHandler(c *gin.Context) error {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

	return nil
}

func (app Application) webhookHandler(c *gin.Context) error {
	println("webhookHandler")
	defer c.Request.Body.Close()

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	var update tgbotapi.Update

	err = json.Unmarshal(bytes, &update)
	if err != nil {
		log.Println(err)
		return err
	}

	// to monitor changes run: heroku logs --tail

	log.Printf("From: %+v Text: %+v\n", update.Message.From, update.Message.Text)
	mc := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	app.Bot.Send(mc)

	return nil
}
