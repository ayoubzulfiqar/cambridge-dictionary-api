package main

import (
	"log"
	"net/http"
)

func main() {
	// Stander JSON REQUEST
	SearchQueryHandler()
	TelegramHandler()
	// Start the HTTP server
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
