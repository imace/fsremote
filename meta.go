package fsremote

type FunMedia struct {
	MediaId     int      `json:"mediaid"`
	Name        string   `json:"name"`
	NameSn      string   `json:"namesn,omitempty"`
	NameEn      string   `json:"nameen,omitempty"`  //name_en
	NameOt      string   `json:"nameot,omitempty"`  //name_ot
	Language    string   `json:"lang,omitempty"`    //language
	MediaLength int      `json:"medialength"`       //medialength
	Country     string   `json:"country,omitempty"` //country
	Release     int64    `json:"release"`           //releasedate
	CoverId     int      `json:"coverid"`
	Tags        []string `json:"tags,omitempty"` //
	Weight      float64  `json:"weight"`         //
	PlayNum     int64    `json:"played"`
}

type FunDoc struct {
	MediaId int      `json:"mediaid"`
	Name    string   `json:"name"`
	Tags    []string `json:"tags,omitempty"`
	Release int64    `json:"release"`
	Weight  int      `json:"weight"`
}
