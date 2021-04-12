package test

import (
	"channel/utils"
	"testing"
)

func TestNotify(t *testing.T) {
	n := utils.GetNotify()
	err := n.SendPopo("test", []string{""})
	if err != nil {
		t.Log(err.Error())
	}

}
