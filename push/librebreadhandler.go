package push

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type LibreBreadHandler struct {
	librePush *LibrePush
}

func NewLibreBreadHandler(librePush *LibrePush) *LibreBreadHandler {
	return &LibreBreadHandler{
		librePush: librePush,
	}
}

func (h *LibreBreadHandler) HandlePush(w http.ResponseWriter, r *http.Request) {
	var msg BatchMessage

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("cannot decode request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var response *BatchResponse

	if msg.ValidateOnly {
		response, err = h.dryRun(msg)
	} else {
		response, err = h.send(msg)
	}

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("cannot encode json reponse: %v", err)
	}
}

func (h *LibreBreadHandler) send(msg BatchMessage) (*BatchResponse, error) {
	response, err := h.librePush.Send(msg)
	if err != nil {
		return nil, fmt.Errorf("cannot emulate push send: %w", err)
	}

	return response, nil
}

func (h *LibreBreadHandler) dryRun(msg BatchMessage) (*BatchResponse, error) {
	response, err := h.librePush.SendDryRun(msg)
	if err != nil {
		return nil, fmt.Errorf("cannot dry run push send: %w", err)
	}

	return response, nil
}
