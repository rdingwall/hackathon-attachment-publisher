package html2png

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	Uri string `json:"uri"`
}

var httpClient = &http.Client{}

func Post(bodyBase64Url string) (*Response, error) {
	postUri := os.Getenv("HTML2PNG_POST_URI")
	log.Printf("POST uri=%s length=%v data=%s\n", postUri, len(bodyBase64Url), bodyBase64Url)

	request, err := http.NewRequest("POST", postUri, strings.NewReader(bodyBase64Url))
	if err != nil {
		return nil, err
	}

	request.Header.Add("content-type", "text/plain")

	res, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(body))
	}

	response := &Response{}
	err = json.NewDecoder(res.Body).Decode(response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
