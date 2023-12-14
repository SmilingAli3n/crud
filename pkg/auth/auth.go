package auth

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/SmilingAli3n/crud/pkg/db"
	_ "github.com/go-sql-driver/mysql"
	bcrypt "golang.org/x/crypto/bcrypt"
)

type User struct {
	id       int64
	Name     string
	Password string
	hash     string
}

func Authorized(r *http.Request) bool {
	username, pass, ok := r.BasicAuth()
	if !ok {
		log.Print(fmt.Errorf("User %s failed to log in", username))
	}

	u := User{
		Name:     username,
		Password: pass,
	}
	if u.exists() && u.passwordIsCorrect() {
		return true
	} else {
		return false
	}
}

func (u *User) exists() bool {
	db, err := db.GetInstance()
	if err != nil {
		return false
	}
	defer db.Close()
	err2 := db.QueryRow("SELECT id FROM users WHERE username = ?", u.Name).Scan(&u.id)
	if err2 == nil {
		return true
	} else if err2 == sql.ErrNoRows {
		return false
	} else {
		return false
	}
}

func (u *User) passwordIsCorrect() bool {
	if u.id == 0 {
		log.Print(fmt.Errorf("User id %d is undefined or incorrect", u.id))
		return false
	}

	db, err := db.GetInstance()
	if err != nil {
		log.Print(err)
		return false
	}
	defer db.Close()
	rows, err := db.Query("SELECT password_hash FROM users WHERE id = ?", u.id)
	if err != nil {
		log.Print(err)
		return false
	}
	if rows.Next() {
		if err := rows.Scan(&u.hash); err != nil {
			log.Print(err)
			return false
		}
	}

	if bcrypt.CompareHashAndPassword([]byte(u.hash), []byte(u.Password)) == nil {
		return true
	} else {
		return false
	}
}
