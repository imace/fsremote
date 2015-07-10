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
}

type FunDoc struct {
	MediaId int      `json:"mediaid"`
	Name    string   `json:"name"`
	Tags    []string `json:"tags,omitempty"`
	Release int64    `json:"release"`
	Weight  int      `json:"weight"`
}

//select mediaid,playnum, daynum,seven_daysnum,weeknum,monthnum,modifydate
type FunTomato struct {
	MediaId  int     `json:"mediaid"`
	Fresh    float64 `json:"fresh"`
	PlayNum  int     `json:"playnum"`
	DayNum   int     `json:"daynum"`
	Day7Num  int     `json:"day7num"`
	WeekNum  int     `json:"weeknum"`
	MonthNum int     `json:"monthnum"`
	Date     int64   `json:"date"`
}
