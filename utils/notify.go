package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const url = ""

//Notify msg
type Notify struct {
	accessKey string
	sender    string
}

type popo struct {
	Type    string   `json:"message_type"`
	Content string   `json:"content"`
	Sender  string   `json:"sender"`
	Reciver []string `json:"reciever_list"`
}

func getDefaultReciever() []string {
	return []string{""}
}

//GetNotify return
func GetNotify() Notify {
	return Notify{accessKey: "", sender: ""}
}

func isExists(source []string, target string) bool {
	for _, v := range source {
		if v == target {
			return true
		}
	}
	return false
}

//SendPopo popo msg
func (notify *Notify) SendPopo(msg string, users []string) (err error) {
	admin := getDefaultReciever()
	for _, v := range admin {
		if !isExists(users, v) {
			users = append(users, v)
		}
	}

	p := &popo{Type: "popo", Content: msg, Sender: notify.sender, Reciver: users}

	b, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	req.Header.Set("X-Notify-AccessKey", notify.accessKey)
	req.Header.Set("Content-Type", "application/json")

	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}
