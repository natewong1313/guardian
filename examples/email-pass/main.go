package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/natewong1313/guardian"
	"github.com/natewong1313/guardian/pkg/adapters/sqlite"
	"golang.org/x/crypto/bcrypt"
)

func setupDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	migrate(db)
	return db, nil
}

type userRow struct {
	ID           int
	Email        string
	PasswordHash string
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER AUTOINCREMENT PRIMARY KEY,
		email TEXT UNIQUE,
		password_hash VARCHAR(64)
	)`)
	if err != nil {
		return err
	}
	return nil
}

type signinRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	adapter, err := sqlite.New("./foo.db")
	if err != nil {
		panic(err)
	}
	db, err := setupDB()
	if err != nil {
		panic(err)
	}

	router := http.NewServeMux()
	router.HandleFunc("POST /signin", func(w http.ResponseWriter, r *http.Request) {
		body := &signinRequest{}
		if err := json.NewDecoder(r.Body).Decode(body); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		user := &userRow{}
		row := db.QueryRow("SELECT * FROM users WHERE email=?", body.Email)
		if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if user.PasswordHash != string(hashedPassword) {
			http.Error(w, "Password is invalid", 401)
			return
		}

		session_token := guardian.GenerateSessionToken()
		session, err := guardian.CreateSession(session_token, user.ID, adapter)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:   "session",
			Value:  session.ID,
			Path:   "/",
			MaxAge: int(session.ExpiresAt.Unix()),
		})

		fmt.Fprintf(w, "success", 201)
	})

	if err := http.ListenAndServe(":6969", router); err != nil {
		panic(err)
	}
}
