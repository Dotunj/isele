package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

var (
	host           string
	endpointID     string
	maskID         string
	apiKey         string
	rate           int
	duration       int
	projectID      string
	body           string
	defaultPayload = json.RawMessage(`{
			"event": "payment.success",
			"data": {
				"status": "Completed",
				"description": "Transaction successful",
				"userID": "test_user_id808",
				"paymentReference": "test_ref_85149",
				"amount": 200,
				"senderAccountName": "Alan Christian Segun",
				"sourceAccountNumber": "299999993564",
				"sourceAccountType": "personal",
				"sourceBankCode": "50211",
				"destinationAccountNumber": "00855584818",
				"destinationBankCode": "063"
			}
		}`)
)

func init() {
	rootCmd.PersistentFlags().StringVar(&host, "host", "http://localhost:5005", "Convoy Base Url")
	rootCmd.PersistentFlags().StringVar(&apiKey, "apiKey", "", "Api Key")
	rootCmd.PersistentFlags().StringVar(&projectID, "projectID", "", "Project ID")
	rootCmd.PersistentFlags().StringVar(&endpointID, "endpointID", "", "Endpoint ID")
	rootCmd.PersistentFlags().StringVar(&maskID, "maskID", "", "Source Mask ID")
	rootCmd.PersistentFlags().IntVar(&rate, "rate", 1, "Number of Events to Create")
	rootCmd.PersistentFlags().IntVar(&duration, "duration", 1, "Duration in seconds")
	rootCmd.PersistentFlags().StringVar(&body, "body", "", "Path to Request JSON file")

}

var rootCmd = &cobra.Command{
	Use:   "serve",
	Short: "Send Events to your Convoy instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		url := fmt.Sprintf("%s/api/v1/projects/%s/events", host, projectID)

		if maskID != "" {
			url = fmt.Sprintf("%s/ingest/%s", host, maskID)
		}

		payload := defaultPayload
		if body != "" {
			jsonFile, err := os.Open(body)
			if err != nil {
				return err
			}

			defer jsonFile.Close()

			payload, err = io.ReadAll(jsonFile)
			if err != nil {
				return err
			}
		}

		fmt.Println("Making request to URL >> ", url)

		// prepare request data
		e := struct {
			EndpointID string          `json:"endpoint_id"`
			Data       json.RawMessage `json:"data"`
			EventType  string          `json:"event_type"`
		}{EndpointID: endpointID, Data: payload, EventType: "*"}

		b, err := json.Marshal(&e)
		if err != nil {
			fmt.Println("error:", err)
			return err
		}

		// number of events to send
		r := vegeta.Rate{Freq: rate, Per: time.Second}
		d := time.Duration(duration) * time.Second
		bearer := fmt.Sprintf("Bearer %s", apiKey)

		h := http.Header{
			"Content-Type":  {"application/json"},
			"Authorization": {bearer},
		}

		targeter := vegeta.NewStaticTargeter(vegeta.Target{
			Method: "POST",
			URL:    url,
			Body:   b,
			Header: h,
		})

		attacker := vegeta.NewAttacker()

		var metrics vegeta.Metrics
		for res := range attacker.Attack(targeter, r, d, "Isele Dispatcher") {
			metrics.Add(res)
		}

		metrics.Close()

		fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)

		return nil
	},
}
