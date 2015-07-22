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

func MediaScore(d, w, s14, m, t int, date int64, typ string) (w1, w2 float64) {
	typw := 1.0
	switch typ {
	case "tv":
		typw = 1.0
	case "movie":
		typw = 1.1
	case "cartoon":
		typw = 0.9
	case "variety":
		typw = 0.8
	case "vfilm":
		typw = 0.7
	}
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
	tw := unit_score(t, 2, 1, 742287)
	mw := unit_score(m, 2, 3, 22605)
	s14w := unit_score(s14, 3, 4, 10063)
	ww := unit_score(w, 5, 6, 5580)
	dw := unit_score(d, 6, 7, 715)
	//台风 33.083 13.344 124 715 10063 5580 22605 742287 1450022400
	x := tw + mw + s14w + ww + dw
	x = math.Log(x)*typw + 1.0
	return x * weight, x
}
func unit_score(v int, C, R float64, m int) float64 {
	return float64(v)/float64(v+m)*R + float64(m)/float64(v+m)*C
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
