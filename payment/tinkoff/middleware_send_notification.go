package tinkoff

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/vasyahuyasa/librebread/payment"
)

const (
	successResponse     = "OK"
	notificationTimeout = 10 * time.Second
)

func sendNotification(paymentProcess payment.PaymentProcess, client payment.Client, provider payment.Provider) (ok bool, err error) {
	notification := map[string]interface{}{
		"Pan":         client.CardNumber,
		"PaymentId":   paymentProcess.ProcessID,
		"Status":      paymentProcess.Status,
		"RebillId":    client.RebillID,
		"Amount":      paymentProcess.Amount,
		"TerminalKey": provider.Login,
	}

	notification["Token"], err = generateToken(notification, provider.Password)
	if err != nil {
		return false, err
	}

	postBody, _ := json.Marshal(notification)
	requestBody := bytes.NewBuffer(postBody)

	httpClient := http.Client{
		Timeout: notificationTimeout,
	}

	resp, err := httpClient.Post(paymentProcess.NotificationURL, "application/json", requestBody)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if string(body) != successResponse {
		return false, nil
	}

	return true, nil
}
