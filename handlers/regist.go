package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"unicode"
	"vtask/database"
	"vtask/internal/user"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var req user.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil{
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if !isValidLogin(req.Login){
		http.Error(w, `{"error": "Login must be 3-30 characters, letters and digits only"}`, http.StatusBadRequest)
		return
	}


	if !isValidPassword(req.Password){
		http.Error(w, `{"error": "Password must be at least 8 characters with 1 uppercase letter and 1 digit"}`, http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil{
		http.Error(w, "Failed to secure password", http.StatusInternalServerError)
		return
	}

	var userID int
	query := `INSERT INTO  vkusers(login, password) VALUES ($1, $2) RETURNING id`
	err = database.DB.QueryRow(query, req.Login, string(hashedPassword)).Scan(&userID)

	   if err != nil {

        log.Printf("DB error: %v", err) 
		
        if pqErr, ok := err.(*pq.Error); ok {
            if pqErr.Code.Name() == "unique_violation" {
                http.Error(w, `{"error": "Login already exists"}`, http.StatusConflict)
                return
            }
        }
        
        http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
        return
    }


	response := user.UserResponse{
		ID:    userID,
		Login: req.Login,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}


func isValidLogin(login string) bool{
	if len(login) < 3 || len(login) > 30{
		return false
	}
	return regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(login)
}

func isValidPassword(password string) bool{
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasDigit := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	return hasUpper && hasDigit

}