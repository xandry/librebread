package push

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type LibreBreadHandler struct {
	librePush *LibrePush
}

func NewLibreBreadHandler(librePush *LibrePush) *LibreBreadHandler {
	return &LibreBreadHandler{
		librePush: librePush,
	}
}

func (h *LibreBreadHandler) HandleSend(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	var msg BatchMessage

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("cannot decode request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.librePush.Send(provider, msg)
	if err != nil {
		log.Printf("cannot emulate push send: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("cannot encode json reponse: %v", err)
	}
}

func (h *LibreBreadHandler) HandleSendDryRun(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	var msg BatchMessage

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("cannot decode request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.librePush.Send(provider, msg)
	if err != nil {
		log.Printf("cannot emulate push send: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("cannot encode json reponse: %v", err)
	}
}
