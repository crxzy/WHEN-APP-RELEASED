package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//CommonRequest external service
type CommonRequest struct {
	Name        string
	BundleID    string
	PackageName string
	Extra       string
}

//CommonResponse common response
type CommonResponse struct {
	Version     string
	ReleaseTime string
	Msg         string
}

//Check the version
func Check(url string, req CommonRequest) (resp CommonResponse, err error) {
	jsonByte, err := json.Marshal(req)
	if err != nil {
		return
	}

	r, err := http.Post(url, "application/json", bytes.NewBuffer(jsonByte))
	if err != nil {
		return
	}

	if r.StatusCode != http.StatusOK {
		err = fmt.Errorf("statusCode:%d", r.StatusCode)
		return
	}

	defer r.Body.Close()
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &resp)
	fmt.Println(resp)

	if resp.Msg != "ok" {
		err = fmt.Errorf(resp.Msg)
	}
	return
}
