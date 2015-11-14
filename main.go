package main

import (
	"github.com/rdingwall/hackathon-attachment-publisher/Godeps/_workspace/src/github.com/go-martini/martini"
	"github.com/rdingwall/hackathon-attachment-publisher/controllers"
	"github.com/rdingwall/hackathon-attachment-publisher/matching"
	"github.com/rdingwall/hackathon-attachment-publisher/mondo"
	"log"
	"github.com/joho/godotenv"
	"os"
)

var matcher = matching.NewMatcher()

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	mondoApiUri := os.Getenv("MONDO_API_URI")
	mondoAccessToken := os.Getenv("MONDO_ACCESS_TOKEN")

	m := martini.Classic()
	mondoApiClient := &mondo.MondoApiClient{Url: mondoApiUri, AccessToken: mondoAccessToken}
	m.Map(matcher)
	m.Map(mondoApiClient)
	m.Post("/webhooks/mondo/transaction", controllers.PostMondoWebhook)
	m.Post("/webhooks/email", controllers.PostEmailWebhook)
	m.Run()
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
