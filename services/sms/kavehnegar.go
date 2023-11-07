package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	kavehnegarApiUrl      = "https://api.kavenegar.com/v1/%s/sms/send.json?receptor=%s&message=%s"
	kavehnegarContentType = "application/x-www-form-urlencoded"
)

type KavehnegarSMSService struct {
	apiKey string
}

func (s *KavehnegarSMSService) SetAPIKey(apikey string) {
	s.apiKey = apikey
}

func (s *KavehnegarSMSService) SendSMS(number string, message string) error {
	reqUrl := fmt.Sprintf(kavehnegarApiUrl, s.apiKey, url.QueryEscape(number), url.QueryEscape(message))
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("content-type", kavehnegarContentType)

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
