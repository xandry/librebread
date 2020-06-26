package helpdesk

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

func HelpdeskEddyHandler(stor *HelpdeskStorage) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msg := HelpdeskMsg{
			Time:         time.Now(),
			TypeId:       atoi(r.FormValue("type_id")),
			PriorityId:   atoi(r.FormValue("priority_id")),
			DepartmentId: atoi(r.FormValue("department_id")),
			Title:        r.FormValue("title"),
			Description:  r.FormValue("description"),
		}

		log.Printf("HelpdeskEddy new: %s %s", msg.Title, msg.Description)

		stor.Push(msg)
	})
}

func atoi(str string) int {
	v, _ := strconv.Atoi(str)
	return v
}
