package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	errorJsonString string = "{ \"status\" : \"-1\"}"
)

func registerUser(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")

	if req.Method != "POST" {
		fmt.Println("Register User Method Not Post")
		json.NewEncoder(res).Encode(errorJsonString)
		return
	}

	errFormParse := req.ParseMultipartForm(32 << 20)

	if errFormParse != nil {
		fmt.Println("Error in Parsing Form")
		json.NewEncoder(res).Encode(errorJsonString)
		return
	}

}

func main() {

	//Serves the Registration form to register the user
	//Routes /Register/ to Views and index.html is invoked
	http.Handle("/User/Register/", http.StripPrefix("/User/Register/", http.FileServer(http.Dir("Views/Register"))))

	//Serves login form to the user
	http.Handle("/User/Login/", http.StripPrefix("/User/Login/", http.FileServer(http.Dir("Views/Login"))))

	http.HandleFunc("/RegisterUser", registerUser)

	http.ListenAndServe(":8008", nil)

}
