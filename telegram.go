package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var TelegramTOKEN string = os.Getenv("TOKEN")

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}
type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}
type Chat struct {
	Id int `json:"id"`
}

func SendTextToTelegramChat(chatId int, text string) (string, error) {
	var telegramApi string = "https://api.telegram.org/bot" + TelegramTOKEN + "/sendMessage"
	parseMode := "Markdown"
	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chatID":    {strconv.Itoa(chatId)},
			"text":      {text},
			"parseMode": {parseMode},
		})

	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = io.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram answer %s", errRead.Error())
		return "", err
	}
	bodyString := string(bodyBytes)

	return bodyString, nil
}
