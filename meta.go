package fsremote

import (
	"math"
	"strconv"
	"time"
)

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
type IntString int
type EsMedia struct {
	Actor       string    `json:"actor"`
	Awards      string    `json:"awards"`
	Country     string    `json:"country"`
	CoverPicID  int       `json:"coverPicId"`
	Day         int       `json:"day"`
	Director    string    `json:"director"`
	DisplayType string    `json:"displayType"`
	FirstCharCn string    `json:"firstCharCn"`
	Image       string    `json:"image"`
	Lang        string    `json:"lang"`
	MediaID     int       `json:"mediaId"`
	MediaLength IntString `json:"mediaLength"`
	Month       int       `json:"month"`
	Name        string    `json:"name"`
	NameEn      string    `json:"nameEn"`
	NameOt      string    `json:"nameOt"`
	PinyinCn    string    `json:"pinyinCn"`
	Play        int       `json:"play"`
	Release     string    `json:"release"`
	Role        string    `json:"role"`
	Seven       int       `json:"seven"`
	Tags        string    `json:"tags"`
	Week        int       `json:"week"`
	Weight      float64   `json:"weight"`
	Weight2     float64   `json:"weight2`
}
type EsMedias struct {
	Data []EsMedia `json:"data"`
	Num  int       `json:"num"`
}

const _2020 = 1577836800

func MediaScore(d, w, s7, m, t int, date int64) float64 {
	r, l, dt := time.Unix(0, 0), time.Unix(_2020, 0), time.Unix(int64(date), 0)
	rg := math.Log(l.Sub(r).Hours() / 24)
	dr := dt.Sub(r).Hours() / 24
	dl := l.Sub(dt).Hours() / 24
	if dr < 1.0 {
		dr = 1.0
	}
	if dl < 1.0 {
		dl = 1.0
	}
	r1, l1 := math.Log(dr)/rg, math.Log(dl)/rg
	weight := r1*1.0 + l1*2.0
	days := time.Since(time.Unix(int64(date), 0)).Hours() / 24
	if days < 1.0 {
		days = 1.0
	}
	x := float64(t)/days + float64(d) + float64(s7)/7.0 + float64(w)/5.0 + float64(m)/30.0
	return x * weight
}
func (m *IntString) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(*m))), nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *IntString) UnmarshalJSON(data []byte) error {
	t, err := strconv.Unquote(string(data))
	if err != nil {
		t = string(data)
	}
	i, err := strconv.Atoi(t)
	*m = IntString(i)
	return err
}
