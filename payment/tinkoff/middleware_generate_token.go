package tinkoff

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
)

// Алгоритм формирования подписи запроса:
// https://www.tinkoff.ru/kassa/develop/api/request-sign/
func generateToken(data map[string]interface{}, password string) (token string, err error) {
	delete(data, "Token")
	delete(data, "Shops")
	delete(data, "Receipt")
	delete(data, "DATA")

	data["Password"] = password

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	strValues := ""
	for _, k := range keys {
		strValues += fmt.Sprint(data[k])
	}

	hash := sha256.New()
	if _, err = hash.Write([]byte(strValues)); err != nil {
		return "", err
	}

	token = hex.EncodeToString(hash.Sum(nil))

	return token, nil
}
