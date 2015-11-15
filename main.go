package main

import (
	"github.com/go-martini/martini"
	"github.com/joho/godotenv"
	"github.com/rdingwall/hackathon-attachment-publisher/controllers"
	"github.com/rdingwall/hackathon-attachment-publisher/matching"
	"github.com/rdingwall/hackathon-attachment-publisher/mondo"
	"log"
	"os"
	"strings"
)

var matcher = matching.NewMatcher()

func main() {
	if martini.Env != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}
	}

	mondoApiUri := os.Getenv("MONDO_API_URI")
	mondoAccessToken := os.Getenv("MONDO_ACCESS_TOKEN")
	addr := os.Getenv("ATTACHMENT_PUBLISHER_ADDR")
	vendors := os.Getenv("ATTACHMENT_PUBLISHER_VENDORS")
	if vendors == "" {
		log.Fatalf("no vendors configured!")
	}
	controllers.Vendors = strings.Split(vendors, ",")

	m := martini.Classic()
	mondoApiClient := &mondo.MondoApiClient{Url: mondoApiUri, AccessToken: mondoAccessToken}
	m.Map(matcher)
	m.Map(mondoApiClient)
	m.Post("/webhooks/mondo/transaction", controllers.PostMondoWebhook)
	m.Post("/webhooks/email", controllers.PostEmailWebhook)
	m.Get("/", func() string {
		return "Hello Mondo crowd!"
	})

	m.RunOnAddr(addr)
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
