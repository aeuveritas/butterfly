package notification

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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

// Notify send notification to telegram
func Notify(filename string, token string) {
	if token != "" {
		baseURL := "https://api.telegram.org/bot" + token + "/"

		initURL := baseURL + "getUpdates"
		resp, err := http.Get(initURL)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		var target initResponse
		err = json.NewDecoder(resp.Body).Decode(&target)
		if err != nil {
			panic(err)
		}
		if len(target.Results) != 1 {
			panic("many results in telegram api")
		}
		chatID := target.Results[0].Message.Chat.ID
		filenameInURL := filename + " is saved"
		sendURL := baseURL + "sendMessage?chat_id=" + strconv.Itoa(chatID) + "&text=" + url.QueryEscape(filenameInURL)
		resp, err = http.Post(sendURL, "application/json", nil)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("telegram is not set")
	}
}
