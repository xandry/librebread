package tinkoff

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	paymentpkg "github.com/vasyahuyasa/librebread/payment"
)

const (
	successResponse     = "OK"
	notificationTimeout = 10 * time.Second
)

func sendNotification(payment paymentpkg.Payment, client paymentpkg.Client, provider paymentpkg.Provider) (ok bool, err error) {
	notification := map[string]interface{}{
		"Pan":         client.CardNumber,
		"PaymentId":   payment.PaymentID,
		"Status":      payment.Status,
		"RebillId":    client.RebillID,
		"Amount":      payment.Amount,
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

	resp, err := httpClient.Post(payment.NotificationURL, "application/json", requestBody)
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
