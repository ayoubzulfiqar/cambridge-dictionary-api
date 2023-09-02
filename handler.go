package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// Stander Search API handler
func SearchQueryHandler() {
	http.HandleFunc("/search/", func(w http.ResponseWriter, r *http.Request) {
		wordToSearch := strings.TrimPrefix(r.URL.Path, "/search/")
		outputJSON := getMeaning(wordToSearch)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encodeError := json.NewEncoder(w).Encode(outputJSON)
		if encodeError != nil {
			panic(encodeError)
		}
	})
}

// Telegram API Handler
func TelegramHandler() {
	// Call out the Telegram API
	http.HandleFunc("/telegram/", func(w http.ResponseWriter, r *http.Request) {
		var body Update
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Println(body.Message.Text)
		outputJSON := getMeaning(body.Message.Text)
		telegramResponse, telegramError := SendTextToTelegramChat(body.Message.Chat.Id, PreprocessingJSONToString(outputJSON))
		if telegramError != nil {
			http.Error(w, telegramError.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		telegramResponseEncodeError := json.NewEncoder(w).Encode(telegramResponse)
		if telegramResponseEncodeError != nil {
			panic(telegramResponseEncodeError)
		}
	})

}
func getMeaning(wordToSearch string) wordMeaning {
	return Search(wordToSearch)
}
