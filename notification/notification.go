package notification

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
)

type initResponse struct {
	Ok      bool         `json:"ok"`
	Results []initResult `json:"result"`
}

type initResult struct {
	Message initMessage `json:"message"`
}

type initMessage struct {
	Chat initChat `json:"chat"`
}

type initChat struct {
	ID int `json:"id,omitempty"`
}

// GetTelegramObject return chat id for telegram
func GetTelegramObject(token string, isDebug bool) (string, *tb.Bot) {
	if token != "" {
		URL := "https://api.telegram.org/bot" + token + "/getUpdates"

		resp, err := http.Get(URL)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		var target initResponse
		err = json.NewDecoder(resp.Body).Decode(&target)
		if err != nil {
			panic(err)
		}

		bot, err := tb.NewBot(tb.Settings{
			URL:   "https://api.telegram.org",
			Token: token,
		})
		if err != nil {
			panic(err)
		}

		if len(target.Results) == 0 {
			fmt.Println("telegram is not ready.")
			return "", nil
		}

		return strconv.Itoa(target.Results[0].Message.Chat.ID), bot
	}
	fmt.Println("telegram is not set.")
	return "", nil
}

// BFRecipient recipient for butterfly
type BFRecipient struct {
	chatID string
}

// Recipient implement for Recipient
func (r BFRecipient) Recipient() string {
	return r.chatID
}

// SendText send notification with text to telegram
func SendText(bot *tb.Bot, chatID string, filename string, isDebug bool) {
	if bot != nil && chatID != "" {
		recipient := BFRecipient{chatID: chatID}

		_, file := filepath.Split(filename)
		message := "start record: " + file
		if !isDebug {
			bot.Send(recipient, message)
		}
	}
}

// SendAudio send notification with audio to telegram
func SendAudio(bot *tb.Bot, chatID string, filename string, isDebug bool) {
	if bot != nil && chatID != "" {
		recipient := BFRecipient{chatID: chatID}

		_, file := filepath.Split(filename)
		audio := &tb.Audio{File: tb.FromDisk(filename), Title: file}
		if !isDebug {
			bot.Send(recipient, audio)
		}
	}
}

// SendVideo send notification with video to telegram
func SendVideo(bot *tb.Bot, chatID string, filename string, isDebug bool) {
	if bot != nil && chatID != "" {
		recipient := BFRecipient{chatID: chatID}

		_, file := filepath.Split(filename)
		video := &tb.Video{File: tb.FromDisk(filename), FileName: file}
		if !isDebug {
			bot.Send(recipient, video)
		}
	}
}
