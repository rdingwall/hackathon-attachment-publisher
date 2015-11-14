package main

import (
	"flag"
	"github.com/go-martini/martini"
	"github.com/rdingwall/mondo-feeditem-publisher/controllers"
	"github.com/rdingwall/mondo-feeditem-publisher/matching"
	"github.com/rdingwall/mondo-feeditem-publisher/mondo"
	"log"
)

var mondoApiUri = flag.String("mondoApiUri", "https://staging-api.gmon.io", "Mondo API URI")
var mondoAccessToken = flag.String("mondoAccessToken", "", "Mondo Access Token")

var matcher = matching.NewMatcher()

func main() {
	flag.Parse()
	if *mondoAccessToken == "" {
		flag.PrintDefaults()
		return
	}

	m := martini.Classic()
	mondoApiClient := &mondo.MondoApiClient{Url: *mondoApiUri, AccessToken: *mondoAccessToken}
	m.Map(matcher)
	m.Map(mondoApiClient)
	m.Post("/webhooks/mondo/transaction", controllers.PostMondoWebhook)
	m.Post("/webhooks/email", controllers.PostEmailWebhook)
	m.Run()
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
