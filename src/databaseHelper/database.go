package databaseHelper

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"model"
)

//Global Database Variable
var Db *sql.DB

func InitDatabase() error {
	var err error
	Db, err = sql.Open("mysql", "root:admin@/gologin")

	if err != nil {
		fmt.Println(err)
		return err
	}

	err = Db.Ping()

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func CheckUserExist(username string) bool {

	query := "SELECT EXISTS(Select 1 from users where username = '" + username + "');"

	var exist bool

	err := Db.QueryRow(query).Scan(&exist)

	if err != nil {
		fmt.Println(err)
		return true
	}
	return exist
}

func AddToDatabase(user model.User) error {

	stmt, err := Db.Prepare("INSERT INTO users SET username = ? , password = ? , first_name = ? , last_name = ? , session_id = ? ")

	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = stmt.Exec(user.Username, user.Password, user.FirstName, user.LastName, user.SessionId)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func CheckPassword(user model.User) error {

	query := "SELECT password FROM users where username = '" + user.Username + "';"
	var password string
	err := Db.QueryRow(query).Scan(&password)

	if err != nil {
		fmt.Println(err)
		return err
	}

	errCompare := bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))

	if errCompare != nil {
		fmt.Println(errCompare)
		return errCompare
	} else {
		fmt.Println("User Authenticated")
		return nil
	}
}

func GetUserDetail(user model.User) (model.User, error) {

	query := "Select first_name, last_name, session_id from users where username = '" + user.Username + "';"

	err := Db.QueryRow(query).Scan(&user.FirstName, &user.LastName, &user.SessionId)

	if err != nil {
		fmt.Println(err)
		return user, err
	}

	return user, nil

}

func AddSessionToDb(sessionId string) error {

	stmt, err := Db.Prepare("Insert into session SET session_id = ?")

	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = stmt.Exec(sessionId)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

func DeleteSessionFromDb(sessionId string) error {

	stmt, err := Db.Prepare("Delete from session where session_id = ?")

	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = stmt.Exec(sessionId)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func CheckSessionExistInDb(sessionId string) bool {

	query := "Select exists( Select 1 from session where session_id = '" + sessionId + "');"

	var exist bool

	err := Db.QueryRow(query).Scan(&exist)

	if err != nil {
		fmt.Println(err)
		return true
	}

	return exist

}
