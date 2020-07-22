package api

import (
	"html/template"
	"log"
	"net/http"
)

type renderer struct {
	smsTpl *template.Template
}

func newRenderer() *renderer {
	r := &renderer{}

	r.smsTpl = parseTemplates("sms", baseTmpl, smsTempl)

	return r
}

func parseTemplates(name string, tmpls ...string) *template.Template {
	t := template.New(name)
	var err error

	for _, s := range tmpls {
		t, err = t.Parse(s)
		if err != nil {
			log.Fatalf("can not parse template %s: %v", name, err)
		}
	}

	return t
}

func (re *renderer) renderSms(w http.ResponseWriter, smses SMSes) error {
	return re.smsTpl.Execute(w, smses)
}
