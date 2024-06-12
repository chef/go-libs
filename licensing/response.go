package licensing

import (
	"log"
	"time"
)

type Response struct {
	Data       interface{} `json:"data"`
	Message    string      `json:"message"`
	StatusCode int         `json:"status_code"`
}

type Client struct {
	License   string `json:"license"`
	Status    string `json:"status"`
	ChangesTo string `json:"changesTo"`
	ChangesOn string `json:"changesOn"`
	ChangesIn int    `json:"changesIn"`
	Usage     string `json:"usage"`
	Used      int    `json:"used"`
	Limit     int    `json:"limit"`
	Measure   string `json:"measure"`
}

type ClientResponse struct {
	Data struct {
		Client Client `json:"client"`
	} `json:"data"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

type ValidateResponse struct {
	Data       bool   `json:"data"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func haveGrace(client Client) bool {
	return client.Status == "Grace"
}

func isExpired(client Client) bool {
	return client.Status == "Expired"
}

func isExhausted(client Client) bool {
	return client.Status == "Exhausted"
}

func isAboutToExpire(client Client) (out bool) {
	expiresOn, err := time.Parse(time.RFC3339, client.ChangesOn)
	if err != nil {
		log.Fatal("Unknown expiration time received from the server: ", err)
	}

	expirationIn := int(time.Until(expiresOn).Hours() / 24)
	return client.Status == "Active" && client.ChangesTo == "Expired" && expirationIn >= 1 && expirationIn <= 7
}

func isExpiringOrExpired(client Client) bool {
	return haveGrace(client) || isExpired(client) || isAboutToExpire(client)
}
