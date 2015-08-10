package main

import (
	"encoding/json"
	"testing"
)

func TestHiProtocol(t *testing.T) {
	var sem = `{
    "code": "SETTING_EXEC",
    "history": "cn.yunzhisheng.setting",
    "text": "储存设置",
    "semantic": {
    "intent": {
    "operator": "ACT_SET"
    }
    },
    "service": "cn.yunzhisheng.setting",
    "rc": 0
    }`
	var s hi_understand_result
	if err := json.Unmarshal([]byte(sem), &s); err != nil {
		t.Error(err)
	}
	t.Log(s)
}
