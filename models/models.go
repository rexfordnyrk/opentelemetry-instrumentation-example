package models

import (
	"io/ioutil"
	"net/http"
	"time"
)

type Parcel struct {
	ID                 string    `gorm:"primaryKey"`
	CustomerID         string    `json:"customer_id"`
	SenderName         string    `json:"sender_name"`
	SenderPhone        string    `json:"sender_phone"`
	OriginAddress      string    `json:"origin_address"`
	OriginCity         string    `json:"origin_city"`
	OriginState        string    `json:"origin_state"`
	DestinationAddress string    `json:"destination_address"`
	DestinationCity    string    `json:"destination_city"`
	DestinationState   string    `json:"destination_state"`
	RecipientName      string    `json:"recipient_name"`
	RecipientPhone     string    `json:"recipient_phone"`
	Weight             int       `json:"weight"`
	Height             int       `json:"height"`
	Width              int       `json:"width"`
	Length             int       `json:"length"`
	Fee                string    `json:"fee"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (p Parcel) GenerateFee() error {
	// API endpoint for generating a random integer between 1 and 100
	apiURL := "https://www.random.org/integers/?num=1&min=1&max=100&col=1&base=10&format=plain&rnd=new"

	// Make the HTTP GET request
	response, err := http.Get(apiURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// Convert the response body to a string and return
	randomNumber := string(body)

	p.Fee = "$" + randomNumber

	return nil
}
