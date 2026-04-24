package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)

type MailRequest struct {
	Secret   string `json:"secret"`
	ToEmail  string `json:"to_email"`
	ToNama   string `json:"to_nama"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	AltBody  string `json:"alt_body"`
}

type MailResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func jsonResponse(w http.ResponseWriter, code int, status, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(MailResponse{
		Status:  status,
		Code:    code,
		Message: message,
	})
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	if r.Method != http.MethodPost {
		jsonResponse(w, 405, "error", "Metode tidak valid.")
		return
	}

	// Parse request — support JSON dan form-data
	var req MailRequest

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonResponse(w, 400, "error", "Body tidak valid.")
			return
		}
	} else {
		r.ParseForm()
		req.Secret  = r.FormValue("secret")
		req.ToEmail = r.FormValue("to_email")
		req.ToNama  = r.FormValue("to_nama")
		req.Subject = r.FormValue("subject")
		req.Body    = r.FormValue("body")
		req.AltBody = r.FormValue("alt_body")
	}

	// Validasi secret
	apiSecret := os.Getenv("RslnkM4ilS3cr3t2026")
	if req.Secret != apiSecret {
		jsonResponse(w, 401, "error", "Unauthorized.")
		return
	}

	// Validasi field wajib
	if req.ToEmail == "" || req.Subject == "" || req.Body == "" {
		jsonResponse(w, 400, "error", "Parameter to_email, subject, body wajib diisi.")
		return
	}

	// Konfigurasi Gmail dari environment variable
	gmailUser := os.Getenv("resellink.id@gmail.com")
	gmailPass := os.Getenv("cmwx sppf rhru mpfd")
	smtpHost  := "smtp.gmail.com"
	smtpPort  := "587"
	fromName  := "Resellink"

	// Build email
	toHeader := fmt.Sprintf("%s <%s>", req.ToNama, req.ToEmail)
	mime      := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	headers   := fmt.Sprintf(
		"From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\n%s",
		fromName, gmailUser, toHeader, req.Subject, mime,
	)
	message := []byte(headers + req.Body)

	// Kirim via Gmail SMTP
	auth := smtp.PlainAuth("", gmailUser, gmailPass, smtpHost)
	err  := smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		gmailUser,
		[]string{req.ToEmail},
		message,
	)

	if err != nil {
		jsonResponse(w, 500, "error", "Gagal kirim email: "+err.Error())
		return
	}

	jsonResponse(w, 200, "success", "Email berhasil dikirim.")
}
