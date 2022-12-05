package api

import (
	"fmt"
	"io"
	"net/http"
)

func MainRoutine() {
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/user", handleUser)
	http.ListenAndServe(":8080", nil)
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("get request")
	io.WriteString(w, "pong")
}
func handleUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getAllUser(w, r)
	case "POST":
		addUser(w, r)
	case "DELETE":
		delUser(w, r)
	case "UPDATE":
		modUser(w, r)
	}
}

//make a logger
//if possible divide client and server

func getAllUser(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Get User")
	response := GetAllUser()
	//pp.Println(response)
	if response.GetErrCode() != 0 {
		http.Error(w, response.GetErrMsg(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response.GetData())
	//w.WriteHeader(200)
}
func addUser(w http.ResponseWriter, r *http.Request) {

}
func delUser(w http.ResponseWriter, r *http.Request) {

}
func modUser(w http.ResponseWriter, r *http.Request) {

}
