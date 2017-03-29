package main

import (
	"databaseHelper"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"model"
	"net/http"
)

var jsonMap = make(map[string]interface{})
var store = sessions.NewCookieStore([]byte("2E9659B26A7E34A3DF672A8BF1613"))

//Helper method to initialise jsonMap with status and message
func initjsonMap(status int, message string) {
	jsonMap["status"] = status
	jsonMap[model.MessageKey] = message
}

/*User Registration*/

//Registers the user in database
//Returns status : 1 for successfull registration
//Returns status : 0 for error or username taken
func registerUser(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")

	if req.Method != "POST" {
		fmt.Println("Register User Method Not Post")
		initjsonMap(0, "Not a Post Request")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	errFormParse := req.ParseMultipartForm(32 << 20)

	if errFormParse != nil {
		fmt.Println("Error in Parsing Form")
		initjsonMap(0, "Error in Prasing Form")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	var newUser model.User

	newUser.FirstName = req.FormValue(model.FirstNameKey)
	newUser.LastName = req.FormValue(model.LastNameKey)
	newUser.Username = req.FormValue(model.UsernameKey)
	newUser.Password = req.FormValue(model.PasswordKey)

	userExist := databaseHelper.CheckUserExist(newUser.Username)

	if userExist {
		fmt.Println("Username already exist")
		initjsonMap(0, "Username already exist")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	pass, errPasswordHash := encryptPassword(&newUser.Password)
	newUser.Password = pass

	if errPasswordHash != nil {
		fmt.Println("Error in creating password hash")
		initjsonMap(0, "Try Again")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	//session-id is created using concatenation of username and firstname
	sess, errSessionIdCreate := createSessionId(newUser.Username, newUser.FirstName)
	newUser.SessionId = sess

	if errSessionIdCreate != nil {
		fmt.Println("Error in creating session-id")
		initjsonMap(0, "Try Again")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	errAddUser := databaseHelper.AddToDatabase(newUser)

	if errAddUser != nil {
		fmt.Println("Error in adding user")
		initjsonMap(0, "Error in adding user")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	initjsonMap(1, "User Registered Successfully")
	json.NewEncoder(res).Encode(jsonMap)
}

//creating session id from username + password concatenation and using bcrypt hash
func createSessionId(username, firstname string) (string, error) {
	sessionId, err := bcrypt.GenerateFromPassword([]byte(username+firstname), bcrypt.MinCost)
	return string(sessionId), err
}

//encrypt password to hash
func encryptPassword(password *string) (string, error) {

	bytePassword, errPasswordHash := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.MinCost)
	return string(bytePassword), errPasswordHash
}

/*User Login*/
func login(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")
	var user model.User

	if req.Method != "POST" {
		fmt.Println("Not a Post Method")
		initjsonMap(0, "Not a Post Method")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	err := req.ParseMultipartForm(32 << 20)

	if err != nil {
		fmt.Println(err)
		initjsonMap(0, "Error in parsing form data")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	user.Username = req.FormValue(model.UsernameKey)
	user.Password = req.FormValue(model.PasswordKey)

	userExist := databaseHelper.CheckUserExist(user.Username)

	if !userExist {
		fmt.Println("User doesnot exist")
		initjsonMap(0, "User doesnot exist")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	err = databaseHelper.CheckPassword(user)

	if err != nil {
		initjsonMap(0, "Password doesnot match")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	user, err = databaseHelper.GetUserDetail(user)

	if err != nil {
		fmt.Println("Error in retreiving file , database error")
		initjsonMap(0, "Error in retreiving file , database error")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	//Checking any other session of the user
	exist := databaseHelper.CheckSessionExistInDb(user.SessionId)

	if exist {
		fmt.Println("A single user can have a single session")
		initjsonMap(0, "A single user can have a single session")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	fmt.Println("User Authenticated")

	session, errSession := store.Get(req, "session")
	if errSession != nil {
		fmt.Println("Session from client cannot be decoded : invalid session ")
		fmt.Println(errSession)
		initjsonMap(0, "Session from client cannot be decoded : invalid session ")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	err = databaseHelper.AddSessionToDb(user.SessionId)

	if err != nil {
		fmt.Println("Error in adding session-id to database")
		initjsonMap(0, "Error in adding session-id to database")
		json.NewEncoder(res).Encode(jsonMap)
		return

	}

	session.Values["session-id"] = user.SessionId
	session.Save(req, res)
	initjsonMap(1, "User Authenticated")
	json.NewEncoder(res).Encode(jsonMap)

}

func main() {

	err := databaseHelper.InitDatabase()
	if err != nil {
		fmt.Println("Database Error")
		return
	}
	defer databaseHelper.Db.Close()

	//Serves the Registration form to register the user
	//Routes /Register/ to Views and index.html is invoked
	http.Handle("/user/register/", http.StripPrefix("/user/register/", http.FileServer(http.Dir("Views/Register"))))
	http.HandleFunc("/registerUser", registerUser)

	//Serves login form to the user
	http.Handle("/user/login/", http.StripPrefix("/user/login/", http.FileServer(http.Dir("Views/Login"))))
	http.HandleFunc("/login", login)

	http.ListenAndServe(":8008", nil)

}
