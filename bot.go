package main

import (
	"fmt"
	"net/http"
	"net/url"
)

var (
	api_key = ""
	api     = "https://api.telegram.org/bot"
	chat_id = ""
)

func init() {
	api_key = Getenv("API_KEY", api_key)
	api = Getenv("API", api)
	chat_id = Getenv("CHAT_ID", chat_id)
}

func SendMessage(msg string) error {
	if chat_id == "" || api_key == "" {
		fmt.Println(msg)
	} else {
		cmd := fmt.Sprintf("%s/sendMessage?chat_id=%s&parse_mode=MarkdownV2&text=%s", api+api_key, chat_id, url.QueryEscape(msg))
		_, err := http.Get(cmd)
		return err
	}

	return nil
}
