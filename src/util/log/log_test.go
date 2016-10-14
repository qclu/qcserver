package log

// These tests are too simple.

import (
	"testing"
	"time"
)

func TestNewLog(t *testing.T) {
	_, err := NewLog("/root/", "gatewayserv", 1)
	if err != nil {
		t.Error(err.Error())
		return
	}
	ll := GetLog()
	//for {
	ll.LogDebug("write:wocao")
	ll.LogError("err:wocao")
	ll.LogWarn("warn:wocao")
	ll.LogInfo("info:haha")
	time.Sleep(time.Second)
	//}

}
