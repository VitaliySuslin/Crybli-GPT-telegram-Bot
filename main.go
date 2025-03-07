package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


const API_URL = "https://free-unoficial-gpt4o-mini-api-g70n.onrender.com/chat/?query="

func sendQueryToAPI(query string) string {
    fullURL := API_URL + url.QueryEscape(query)

    response, err := http.Get(fullURL)
    if err != nil {
        return "An error occurred while accessing the server."
    }
    defer response.Body.Close()

    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return "ПAn error occurred while reading the response from the server."
    }

    var result map[string]string
    err = json.Unmarshal(body, &result)
    if err != nil {
        return "An error occurred while processing the response from the server."
    }

    return result["results"]
}

func main() {
    botToken := "YOUR_TELEGRAM_BOT_TOKEN"

    bot, err := tgbotapi.NewBotAPI(botToken)
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = true

    log.Printf("Authorized as %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message == nil {
            continue
        }

        userMessage := update.Message.Text

        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please wait while I process your request...")
        bot.Send(msg)

        botResponse := sendQueryToAPI(userMessage)

        msg = tgbotapi.NewMessage(update.Message.Chat.ID, botResponse)
        bot.Send(msg)
    }
}