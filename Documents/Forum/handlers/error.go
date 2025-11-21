package forumino

import (
	"bytes"
	"errors"
	"net/http"
	"syscall"
	"text/template"

	forumino "forumino/models"
)

func RenderError(w http.ResponseWriter, messege string, code int) {
	w.WriteHeader(code)
	temp, err := template.ParseFiles("template/error.html")
	if err != nil {
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		return
	}
	Data := forumino.ErrorPage{
		Messege: messege,
		Code:    code,
	}
	var buff bytes.Buffer

	if err = temp.Execute(&buff, Data); err != nil {
		http.Error(w, "Template execute error", http.StatusInternalServerError)
		return
	}
	if _, err = buff.WriteTo(w); err != nil {
		if !errors.Is(err, syscall.EPIPE) {
			http.Error(w, "FAILED TO WRITE BUFFER", http.StatusInternalServerError)
			return
		}
	}
}
