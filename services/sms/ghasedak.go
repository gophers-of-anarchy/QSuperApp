package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	GhasedakApiUrl       = "https://api.ghasedak.me/v2/sms/send/simple"
	GhasedakContentType  = "application/x-www-form-urlencoded"
	GhasedakCacheControl = "no-cache"
	GhasedakLineNumber   = "30005088"
)

type GhasedakSMSService struct {
	apiKey string
}

func (s *GhasedakSMSService) SetAPIKey(apikey string) {
	s.apiKey = apikey
}

func (s *GhasedakSMSService) SendSMS(number string, message string) error {
	data := url.Values{}
	data.Set("message", message)
	data.Set("receptor", number)
	data.Set("linenumber", GhasedakLineNumber)

	req, err := http.NewRequest(http.MethodPost, GhasedakApiUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("content-type", GhasedakContentType)
	req.Header.Add("apikey", s.apiKey)
	req.Header.Add("cache-control", GhasedakCacheControl)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	log.Println(string(body))

	return nil
}
