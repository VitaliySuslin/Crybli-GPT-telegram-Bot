package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "os"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


const API_URL = "https://free-unoficial-gpt4o-mini-api-g70n.onrender.com/chat/?query="

func sendQueryToAPI(query string) string {
    fullURL := API_URL + url.QueryEscape(query)
    log.Printf("Sending request to API: %s", fullURL)

    response, err := http.Get(fullURL)
    if err != nil {
        log.Printf("Error accessing API: %v", err)
        return "An error occurred while accessing the server."
    }
    defer response.Body.Close()

    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Printf("Error reading response body: %v", err)
        return "An error occurred while reading the response from the server."
    }

    var result map[string]string
    err = json.Unmarshal(body, &result)
    if err != nil {
        log.Printf("Error unmarshaling JSON: %v", err)
        return "An error occurred while processing the response from the server."
    }

    log.Printf("API response received successfully")
    return result["results"]
}

func main() {
    botToken := os.Getenv("BOT_TOKEN")

    if botToken == "" {
        log.Fatal("BOT_TOKEN environment variable is not set")
    }

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
        log.Printf("Received message from user %d: %s", update.Message.From.ID, userMessage)

        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please wait while I process your request...")
        if _, err := bot.Send(msg); err != nil {
            log.Printf("Error sending waiting message: %v", err)
        }

        botResponse := sendQueryToAPI(userMessage)
        log.Printf("Sending response to user %d: %s", update.Message.From.ID, botResponse)

        msg = tgbotapi.NewMessage(update.Message.Chat.ID, botResponse)
        if _, err := bot.Send(msg); err != nil {
            log.Printf("Error sending response message: %v", err)
        }
    }
}