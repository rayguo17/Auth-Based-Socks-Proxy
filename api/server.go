package api

import (
	"encoding/json"
	"fmt"
	"github.com/rayguo17/go-socks/manager/user"
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
//if possible divide client and light

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
	var u user.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res := AddUser(&u)
	if res.GetErrCode() != 0 {
		http.Error(w, res.GetErrMsg(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res.GetData())
}
func delUser(w http.ResponseWriter, r *http.Request) {
	var d DelParams
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//pp.Println(d)
	res := DelUser(&d)
	if res.GetErrCode() != 0 {
		http.Error(w, res.GetErrMsg(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func modUser(w http.ResponseWriter, r *http.Request) {

}
