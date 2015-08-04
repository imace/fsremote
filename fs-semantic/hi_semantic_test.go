package main

import (
	"encoding/json"
	"testing"
)

func TestHiSemantic(t *testing.T) {
	var sem = `{
    "code": "SETTING_EXEC_TV",
    "history": "cn.yunzhisheng.setting.tv",
    "text": "看中央1",
    "semantic": {
    "intent": {
    "operator": "ACT_OPEN_CHANNEL",
    "operands": "CCTV-1",
    "value": "中央一"
    }
    },
    "service": "cn.yunzhisheng.setting.tv",
    "rc": 0
    }`
	var s hi_understand_result
	if err := json.Unmarshal([]byte(sem), &s); err != nil {
		t.Error(err)
	}
	t.Log(s.Semantic.Intent)
}
