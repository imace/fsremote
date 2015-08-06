package main

type hi_protocol struct {
	Domain string                 `json:"domain"`
	Type   string                 `json:"type"`
	Data   map[string]interface{} `json:"data"`
}

//WAITING
type hi_waiting struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	OnCancel string `json:"onCancel"` //[{"domain":"DOMAIN_LOCAL","confirm":"cancel","message":"取消"}]
}
type hi_on_cancel struct {
	Domain  string `json:"domain"`
	Confirm string `json:"confirm"`
	Message string `json:"message"`
}
