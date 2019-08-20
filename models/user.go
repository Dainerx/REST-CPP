package models

import (
	"database/sql"

	"github.com/Dainerx/rest-go-cpp/pkg/password"
)

type User struct {
	Id       int64
	Email    string
	Password string
}

func AllUsers(db *sql.DB) ([]*User, error) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		user := new(User)
		err := rows.Scan(&user.Id, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func AddUser(db *sql.DB, email, pass string) error {
	hash, err := password.HashPass(pass)
	if err != nil {
		return err
	} else {
		_, err := db.Exec("INSERT INTO users (email,password) VALUES(?,?)", email, hash)
		if err != nil {
			return err
		}
		return nil
	}
}

func GetUser(db *sql.DB, id int64) (User, error) {
	rows, err := db.Query("SELECT * FROM users where id=?", id)
	user := new(User)
	if err != nil {
		return *user, err
	} else {
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&user.Id, &user.Email, &user.Password)
			if err != nil {
				return *user, err
			}
		}
		return *user, nil
	}
}

//remove bool as return type leave only user since it reflects the result
func UserExists(db *sql.DB, email, pass string) (bool, *User, error) {
	user := new(User)
	rows, err := db.Query("SELECT id, email, password FROM users where email=?", email)
	if err != nil {
		return false, user, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Email, &user.Password)
		if err != nil {
			return false, user, err
		}
	}
	if ok := password.CheckPassHash(pass, user.Password); !ok {
		return false, user, nil
	}
	return true, user, nil
}
