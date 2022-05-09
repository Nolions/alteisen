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
	e.POST("/alert", ErrHandler(app.alertHandler))
}

func (app Application) healthHandler(c *gin.Context) error {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

	return nil
}

func (app Application) webhookHandler(c *gin.Context) error {
	log.Printf("webhookHandler\n")
	defer c.Request.Body.Close()

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("error:%v\n", err)
		return err
	}

	var update tgbotapi.Update

	err = json.Unmarshal(bytes, &update)
	if err != nil {
		log.Printf("json decode error:%v\n", err)
		return err
	}

	// to monitor changes run: heroku logs --tail
	log.Printf("From: %+v Text: %+v\n", update.Message.From, update.Message.Text)

	return nil
}

type alertReq struct {
	ChatId int64  `json:"chat_id,omitempty"`
	Msg    string `json:"msg"`
	//Service string `json:"service"`
}

func (app Application) alertHandler(c *gin.Context) error {
	log.Printf("alert\n")
	defer c.Request.Body.Close()

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("error:%v\n", err)
		return err
	}

	var a alertReq
	err = json.Unmarshal(bytes, &a)
	if err != nil {
		log.Printf("json decode error:%v\n", err)
		return err
	}

	if a.ChatId == 0 {
		a.ChatId = app.Conf.TargetChatId
	}

	mc := tgbotapi.NewMessage(a.ChatId, a.Msg)
	_, err = app.Bot.Send(mc)
	if err != nil {
		log.Printf("send message error:%v\n", err)
		return err
	}

	return nil
}
