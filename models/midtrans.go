package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type JSON map[string]interface{}

type Response struct {
	PaymentType       string `json:"payment_type"`
	TransactionStatus string `json:"transaction_status"`
}

func RequestPayment(OrderId int, amount int) (redirectURL string, err error) {

	if OrderId <= 0 || amount <= 0 {
		err := errors.New("negative or 0 number")
		return "", err
	}

	url := "https://app.sandbox.midtrans.com/snap/v1/transactions"
	method := "POST"

	payload, err := json.Marshal(JSON{
		"transaction_details": JSON{"order_id": fmt.Sprint("pay-", OrderId), "gross_amount": amount},
		"credit_card":         JSON{"secure": true},
	})
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic U0ItTWlkLXNlcnZlci1LTjFYOFVBckRIdERvcEc1aVF3d1c2Zi0=")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println(string(body))
	temp2 := string(body)
	temp1 := strings.Index(temp2, "https")
	return string(body[temp1 : len(body)-2]), nil
}

func StatusPayment(OrderId string) (redirectURL string, resp Response, err error) {
	url := "https://api.sandbox.midtrans.com/v2/" + OrderId + "/status"
	method := "GET"

	payload := strings.NewReader("\n\n")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic U0ItTWlkLXNlcnZlci1LTjFYOFVBckRIdERvcEc1aVF3d1c2Zi0=")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var response Response
	json.Unmarshal(body, &response)
	fmt.Println(response)

	return string(body), response, nil
}
