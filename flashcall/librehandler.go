package flashcall

import (
	"net/http"
)

func LibrecallHandler(librecall *Librecall) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		strPhone := r.FormValue("phone")
		strCode := r.FormValue("code")

		phone, err := newPhone(strPhone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		code, err := newCode(strCode)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		librecall.Call(phone, code)
	})
}
