package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/rdingwall/hackathon-attachment-publisher/html2png"
	"github.com/rdingwall/hackathon-attachment-publisher/matching"
	"github.com/rdingwall/hackathon-attachment-publisher/mondo"
	"log"
	"net/http"
	"strings"
"github.com/rdingwall/hackathon-attachment-publisher/email"
	"encoding/base64"
)

var Vendors []string = nil

func PostMondoWebhook(w http.ResponseWriter, r *http.Request, matcher *matching.Matcher, mondoApiClient *mondo.MondoApiClient) {
	defer r.Body.Close()
	var request = &mondo.WebhookRequest{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("json parse error: %s\n", err.Error())
		return
	}

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
	if match == nil {
		return
	}

	log.Printf("got match!!")
	err = uploadAttachment(mondoApiClient, match.Transaction.Id, match.Email.BodyHtmlBase64Url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("upload error: %s\n", err.Error())
		return
	}
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

	vendorMatchKey := getVendorMatchKey(request.Subject)
	if vendorMatchKey == "" {
		log.Printf("ignored unrecognised vendor")
		return
	}

	bodyHtml, err := request.HTML()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("html decode error: %s\n", err.Error())
		return
	}

	transaction := &matching.Email{
		MessageId:         request.Id,
		Subject:           request.Subject,
		BodyHtml:          bodyHtml,
		BodyHtmlBase64Url: base64.URLEncoding.EncodeToString([]byte(bodyHtml)),
		VendorMatchKey:    vendorMatchKey,
	}

	match := matcher.MatchEmail(transaction)
	if match == nil {
		return
	}

	err = uploadAttachment(mondoApiClient, match.Transaction.Id, match.Email.BodyHtmlBase64Url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("upload error: %s\n", err.Error())
		return
	}
}

func uploadAttachment(mondoApiClient *mondo.MondoApiClient, transactionId, bodyHtmlBase64 string) error {
	response, err := html2png.Post(bodyHtmlBase64)
	if err != nil {
		return err
	}

	resp, err := mondoApiClient.RegisterAttachment(transactionId, response.Uri, "image/png")
	log.Printf("%v\n", resp)
	return err
}

func formatAmount(amount int32) string {
	a := float32(-amount) / 100
	return fmt.Sprintf("%.2f", a)
}

func getVendorMatchKey(str string) string {
	s := strings.ToLower(str)
	for _, vendorMatchKey := range Vendors {
		if vendorMatchKey != "" && strings.Contains(s, vendorMatchKey) {
			return vendorMatchKey
		}
	}

	return ""
}
