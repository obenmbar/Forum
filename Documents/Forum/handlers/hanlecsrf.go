package forumino

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)


func GenerateCSRFToken() string {
	Bytes := make([]byte, 32)
	rand.Read(Bytes)
	return hex.EncodeToString(Bytes)
}


func SetCSRFToken(w http.ResponseWriter) string {
	token := GenerateCSRFToken()


	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    token,
		HttpOnly: true,   
		Path:     "/",         
		MaxAge:   3600 * 24, 
	
	})

	return token
}


func GetCSRFToken(r *http.Request) string {
	cookie, err := r.Cookie("csrf_token")
	if err != nil {
		return ""
	}
	return cookie.Value
}
func ValidateCSRF(r *http.Request) bool {

	cookie, err := r.Cookie("csrf_token")
	if err != nil {
		return false 
	}

	
	formToken := r.FormValue("csrf_token")


	if cookie.Value == "" || formToken == "" || cookie.Value != formToken {
		return false 
	}

	return true 
}