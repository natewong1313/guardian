package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	if err := migrate(db); err != nil {
		return nil, err
	}
	return db, nil
}

type userRow struct {
	ID           int
	Email        string
	PasswordHash string
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER INCREMENT PRIMARY KEY,
		email TEXT UNIQUE,
		password_hash VARCHAR(64)
	)`)
	if err != nil {
		return err
	}
	return nil
}

type requestBody struct {
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
	router.HandleFunc("GET /protected", func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()

		var sessionToken string
		for _, cookie := range cookies {
			if cookie.Name == "session" {
				sessionToken = cookie.Value
				break
			}
		}
		if sessionToken == "" {
			http.Error(w, "missing session", 401)
			return
		}

		session, err := guardian.ValidateSessionToken(sessionToken, adapter)
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		fmt.Fprintf(w, "user id: %s", session.UserID)
	})
	router.HandleFunc("POST /signup", func(w http.ResponseWriter, r *http.Request) {
		body := &requestBody{}
		if err := json.NewDecoder(r.Body).Decode(body); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		res, err := db.Exec("INSERT INTO users (email, password_hash) VALUES (?, ?)", body.Email, hashedPassword)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		userID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		sessionToken := guardian.GenerateSessionToken()
		session, err := guardian.CreateSession(sessionToken, strconv.Itoa(int(userID)), adapter)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:   "session",
			Value:  sessionToken,
			Path:   "/",
			MaxAge: int(session.ExpiresAt.Unix()),
		})

		w.WriteHeader(201)
		fmt.Fprint(w, "success")
	})
	router.HandleFunc("POST /signin", func(w http.ResponseWriter, r *http.Request) {
		body := &requestBody{}
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

		sessionToken := guardian.GenerateSessionToken()
		session, err := guardian.CreateSession(sessionToken, strconv.Itoa(user.ID), adapter)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:   "session",
			Value:  string(sessionToken),
			Path:   "/",
			MaxAge: int(session.ExpiresAt.Unix()),
		})

		fmt.Fprint(w, "success")
	})

	if err := http.ListenAndServe(":6969", router); err != nil {
		panic(err)
	}
}
