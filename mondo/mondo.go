package mondo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	Authorization   = "Authorization"
	Bearer          = "Bearer "
	ContentType     = "Content-Type"
	ApplicationJson = "application/json"
)

type MondoApiClient struct {
	ClientId     string
	ClientSecret string
	Url          string
	AccessToken  string
}

type RegisterWebhookRequest struct {
	AccountId string `json:"account_id"`
	Url       string `json:"url"`
}

type Webhook struct {
	AccountId string `json:"account_id"`
	Id        string `json:"id"`
	Url       string `json:"url"`
}

type RegisterWebhookResponse struct {
	Webhook Webhook `json:"webhook"`
}

type WebhookRequest struct {
	Type string      `json:"type"`
	Data WebhookData `json:"data"`
}

type WebhookData struct {
	Amount      int32  `json:"amount"`
	Created     string `json:"created"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
	Id          string `json:"id"`
}

type RegisterAttachmentResponse struct {
	Attachment Attachment
}

type Attachment struct {
	Id         string `json:"id"`
	UserId     string `json:"user_id"`
	ExternalId string `json:"external_id"`
	FileUrl    string `json:"file_url"`
	FileType   string `json:"file_type"`
	Created    string `json:"created"`
}

var httpClient = &http.Client{}

func (c *MondoApiClient) RegisterWebHook(accessToken, accountId, webhookUrl string) (*RegisterWebhookResponse, error) {
	log.Printf("Registering webhook for accountId=%s accessToken=%s url=%s\n", accountId, accessToken, webhookUrl)

	webhooksUrl := fmt.Sprintf("%s/webhooks", c.Url)
	formValues := url.Values{
		"account_id": {accountId},
		"url":        {webhookUrl},
	}

	request, err := http.NewRequest("POST", webhooksUrl, strings.NewReader(formValues.Encode()))
	if err != nil {
		return nil, err
	}

	request.Header.Add(Authorization, Bearer+accessToken)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(body))
	}

	webhookResponse := &RegisterWebhookResponse{}
	err = json.NewDecoder(response.Body).Decode(webhookResponse)
	if err != nil {
		return nil, err
	}

	return webhookResponse, nil
}

func (c *MondoApiClient) UnregisterWebHook(accessToken, webhookId string) error {
	log.Printf("Unregistering webhook accessToken=%s webhookId=%s\n", accessToken, webhookId)

	webhooksUrl := fmt.Sprintf("%s/webhooks/%s", c.Url, webhookId)

	request, err := http.NewRequest("DELETE", webhooksUrl, nil)
	if err != nil {
		return err
	}

	request.Header.Add(Authorization, Bearer+accessToken)

	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}

	return nil
}

func (c *MondoApiClient) CreateFeedItem(accountId, itemType, title, imageUrl, body string) error {
	log.Printf("Creating feed item for accountId=%s type=%s title=%s imageUrl=%s body=%s\n", accountId, itemType, title, imageUrl, body)

	feedUrl := fmt.Sprintf("%s/feed", c.Url)
	formValues := url.Values{
		"account_id": {accountId},
		"type":       {itemType},
		"title":      {title},
		"image_url":  {imageUrl},
		"body":       {body},
	}

	request, err := http.NewRequest("POST", feedUrl, strings.NewReader(formValues.Encode()))
	if err != nil {
		return err
	}

	request.Header.Add(Authorization, Bearer+c.AccessToken)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}

	return err
}

func (c *MondoApiClient) RegisterAttachment(externalId, fileUrl, fileType string) (*RegisterAttachmentResponse, error) {
	log.Printf("Registering attachment for access token=%s externalId=%s fileUrl=%s fileType=%s\n", c.AccessToken, externalId, fileUrl, fileType)

	feedUrl := fmt.Sprintf("%s/attachment/register", c.Url)
	formValues := url.Values{
		"external_id": {externalId},
		"file_type":   {fileType},
		"file_url":    {fileUrl},
	}

	request, err := http.NewRequest("POST", feedUrl, strings.NewReader(formValues.Encode()))
	if err != nil {
		return nil, err
	}

	request.Header.Add(Authorization, Bearer+c.AccessToken)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(body))
	}

	registerAttachmentResponse := &RegisterAttachmentResponse{}
	err = json.NewDecoder(response.Body).Decode(registerAttachmentResponse)
	if err != nil {
		return nil, err
	}

	return registerAttachmentResponse, nil
}
