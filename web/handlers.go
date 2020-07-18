package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/vasyahuyasa/librebread/sms"
)

const (
	defaultLimit = 50
)

func redirect(url string, code int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, code)
	})
}

func smsIndexHandler(s *sms.SqliteStorage, re *renderer) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var lim int64 = defaultLimit

		strLim := r.FormValue("limit")
		if strLim != "" {
			var err error

			lim, err = strconv.ParseInt(strLim, 10, 64)
			if err != nil {
				http.Error(w, fmt.Sprintf("can not parse limit: %v", err), http.StatusBadRequest)
				log.Printf("can not parse limit: %v", err)
				return
			}
		}

		messages, err := s.LastMessages(lim)
		if err != nil {
			http.Error(w, fmt.Sprintf("can not get messages: %v", err), http.StatusInternalServerError)
			log.Printf("can not get messages: %v", err)
			return
		}

		if isJson(r) {
			enc := json.NewEncoder(w)

			err = enc.Encode(messages)
			if err != nil {
				http.Error(w, fmt.Sprintf("can not encode messages: %v", err), http.StatusInternalServerError)
				log.Printf("can not encode messages: %v", err)
			}

			return
		}

		err = re.renderSms(w, messages)
		if err != nil {
			http.Error(w, fmt.Sprintf("can not render messages: %v", err), http.StatusInternalServerError)
			log.Printf("can not render messages: %v", err)
		}
	})
}

func isJson(r *http.Request) bool {
	str := strings.ToLower(r.FormValue("json"))

	return str == "1" || str == "true"
}
