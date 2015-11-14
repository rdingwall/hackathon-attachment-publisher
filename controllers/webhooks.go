package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/rdingwall/hackathon-attachment-publisher/email"
	"github.com/rdingwall/hackathon-attachment-publisher/matching"
	"github.com/rdingwall/hackathon-attachment-publisher/mondo"
	"log"
	"net/http"
	"strings"
)

var vendors []string = []string{"beatport", "amazon", "apple"}

func PostMondoWebhook(w http.ResponseWriter, r *http.Request, matcher *matching.Matcher, mondoApiClient *mondo.MondoApiClient) {
	defer r.Body.Close()
	var request = &mondo.WebhookRequest{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("json parse error: %s\n", err.Error())
		return
	}

	go func() {
		if request.Data.Amount > 0 {
			log.Printf("ignored credit transaction")
			return
		}

		vendorMatchKey := getVendorMatchKey(request.Data.Description)
		if vendorMatchKey == "" {
			log.Printf("ignored unrecognised vendor")
			return
		}

		transaction := &matching.Transaction{
			Amount:         formatAmount(request.Data.Amount),
			Created:        request.Data.Created,
			Currency:       request.Data.Currency,
			Description:    request.Data.Description,
			Id:             request.Data.Id,
			VendorMatchKey: vendorMatchKey,
		}

		match := matcher.MatchTransaction(transaction)
		if match != nil {
			log.Printf("got match!!")
			resp, err := mondoApiClient.RegisterAttachment(match.Transaction.Id, "http://i.imgur.com/OLEsqBH.png", "image/png")
			if (err != nil) {
				log.Printf(err.Error())
			}

			log.Printf("%s\n", resp)
		}
	}()
}

func PostEmailWebhook(w http.ResponseWriter, r *http.Request, matcher *matching.Matcher, mondoApiClient *mondo.MondoApiClient) {
	defer r.Body.Close()
	var request = &email.Email{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("json parse error: %s\n", err.Error())
		return
	}

	go func() {

		vendorMatchKey := getVendorMatchKey(request.Subject)
		if vendorMatchKey == "" {
			log.Printf("ignored unrecognised vendor")
			return
		}

		bodyHtml, err := request.HTML()
		if err != nil {
			log.Printf("html decode error: %s\n", err.Error())
			return
		}

		transaction := &matching.Email{
			MessageId:      request.Id,
			Subject:        request.Subject,
			BodyHtml:       bodyHtml,
			VendorMatchKey: vendorMatchKey,
		}

		match := matcher.MatchEmail(transaction)
		if match != nil {
			log.Printf("got match!!")
			resp, err := mondoApiClient.RegisterAttachment(match.Transaction.Id, "http://i.imgur.com/OLEsqBH.png", "image/png")
			if (err != nil) {
				log.Printf(err.Error())
			}

			log.Printf("%s\n", resp)
		}
	}()
}

func formatAmount(amount int32) string {
	a := float32(-amount) / 100
	return fmt.Sprintf("£%.2f", a)
}

func getVendorMatchKey(str string) string {
	s := strings.ToLower(str)
	for _, vendorMatchKey := range vendors {
		if strings.Contains(s, vendorMatchKey) {
			return vendorMatchKey
		}
	}

	return ""
}
