package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"math/rand"
	"net/http"
	"os"
	"parcels/models"
	"sync"
	"time"
)

const (
	apiURL       = "http://localhost:8080/parcels" // Change this to your API URL
	jsonFilePath = "parcels.json"                  // Change this to your JSON file path
)

func main() {
	// Read JSON data from the file
	jsonData, err := os.ReadFile("parcels.json")
	if err != nil {
		panic("Failed to read JSON data")
	}

	// Unmarshal JSON data into a slice of parcels
	var parcels []models.Parcel
	err = json.Unmarshal(jsonData, &parcels)
	if err != nil {
		panic("Failed to unmarshal JSON data")
	}

	// Process parcels in batches of 5
	batchSize := 2
	for i := 0; i < len(parcels); i += batchSize {
		end := i + batchSize
		if end > len(parcels) {
			end = len(parcels)
		}

		// Create a channel to synchronize waiting for the 10-second interval
		waitForInterval := make(chan struct{})

		// Start processing the batch concurrently
		go processBatch(parcels[i:end], waitForInterval)

		// Wait for 10 seconds before starting the next batch
		time.Sleep(1 * time.Second)
		waitForInterval <- struct{}{}
	}
}

func processBatch(batch []models.Parcel, waitForInterval <-chan struct{}) {
	var waitGroup sync.WaitGroup
	for _, parcel := range batch {

		// Start processing the parcel concurrently
		waitGroup.Add(1)
		go func(p models.Parcel) {

			defer waitGroup.Done()

			// POST request
			postParcel(&p)

			// Wait for the 3-second interval to complete
			time.Sleep(2 * time.Second)
			// GET request
			getParcel(p.ID)

			// PUT request with modified values
			modifyParcel(&p)
			putParcel(&p)
		}(parcel)
	}

	// Wait for the 10-second interval to complete
	<-waitForInterval
}

func postParcel(parcel *models.Parcel) {
	//parcel.ID = uuid.New().String()
	jsonData, err := json.Marshal(parcel)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("POST Response Status:", resp.Status)
}

func getParcel(id string) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", apiURL, id))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("GET Response Status:", resp.Status)
}

func modifyParcel(parcel *models.Parcel) {
	// Modify the weight, height, width, and length to random values between 5 and 50
	parcel.Weight = rand.Intn(46) + 5
	parcel.Height = rand.Intn(46) + 5
	parcel.Width = rand.Intn(46) + 5
	parcel.Length = rand.Intn(46) + 5
}

func putParcel(parcel *models.Parcel) {
	jsonData, err := json.Marshal(structs.Map(parcel))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	client := http.Client{}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", apiURL, parcel.ID), bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("PUT Response Status:", resp.Status)
}
