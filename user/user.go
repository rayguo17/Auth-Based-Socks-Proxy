package user

import (
	"encoding/json"
	"github.com/k0kubun/pp/v3"
	"log"
	"os"
)

//upper case so json pkg have access to field
type User struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	CmdChannel chan string
}

var Users []*User
var filePath string = "./user.json"

func ListUser() {
	for i, v := range Users {
		pp.Printf("%d: %v\n", i, v)
	}
}

func init() {
	//read from json file, then form user group
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("read file failed")
		return
	}
	err = json.Unmarshal(fileBytes, &Users)
	if err != nil {
		log.Fatal("unmarshal failed", err.Error())
	}
	pp.Println(Users)
}
