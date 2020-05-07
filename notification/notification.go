package notification

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/aeuveritas/butterfly/parser"
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
	Title string `json:"title,omitempty"`
	ID    int    `json:"id,omitempty"`
}

// GetTelegramBot get telegram bot
func GetTelegramBot(token string, isDebug bool) *tb.Bot {
	if token != "" {
		bot, err := tb.NewBot(tb.Settings{
			URL:   "https://api.telegram.org",
			Token: token,
		})
		if err != nil {
			panic(err)
		}
		return bot
	}

	return nil
}

// GetChatID get chat ID
func GetChatID(tData *parser.TelegramInfo, pData parser.PresetInfo) string {
	if tData.Token != "" {
		targetTitle := pData.Title

		for _, chat := range tData.Chats {
			if targetTitle == chat.Title {
				return chat.ID
			}
		}

		URL := "https://api.telegram.org/bot" + tData.Token + "/getUpdates"

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

		if len(target.Results) == 0 {
			fmt.Println("telegram is not ready.")
			return ""
		}

		for _, result := range target.Results {
			if result.Message.Chat.Title == targetTitle {
				chatID := strconv.Itoa(result.Message.Chat.ID)
				tData.AddItem(targetTitle, chatID)

				return strconv.Itoa(result.Message.Chat.ID)
			}
		}

		fmt.Println("no title in current telegram context")
		return ""
	}

	return ""
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
