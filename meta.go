package fsremote

import "strconv"

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
	Release     IntString `json:"release"`
	Role        string    `json:"role"`
	Seven       int       `json:"seven"`
	Tags        string    `json:"tags"`
	Week        int       `json:"week"`
	Weight      float64   `json:"weight"`
	Weight2     float64   `json:"weight2`
	Score       float64   `json:"_score"`
	Actors      []string  `json:"actors,omitempty"`
	Roles       []string  `json:"roles,omitempty"`
}
type EsMedias struct {
	Data []EsMedia `json:"data"`
	Num  int       `json:"num"`
}

//1437614323
const _1990 = 631152000

func MediaScore(d, w, s14, m, t int, date int64, typ string) (w1, w2 float64) {
	typw := 1.0
	switch typ {
	case "tv":
		typw = 1.0
	case "movie":
		typw = 1.03
	case "cartoon":
		typw = 0.96
	case "variety":
		typw = 0.92
	case "vfilm":
		typw = 0.7
	}

	weight := 1.0 + float64(date-_1990)/60/60/24/11000
	if weight <= 0 {
		weight = 1.0
	}

	tw := unit_score(t, 1, 2, 742287)
	mw := unit_score(m, 1, 3, 22605)
	s14w := unit_score(s14, 1, 5, 10063)
	ww := unit_score(w, 1, 7, 5580)
	dw := unit_score(d, 1, 9, 715)
	//台风 33.083 13.344 124 715 10063 5580 22605 742287 1450022400
	x := tw + mw + s14w + ww + dw
	return x * weight * typw, x * typw
}
func unit_score(v int, C, R float64, m int) float64 {
	return R*float64(v)/float64(v+m) + C*float64(m)/float64(v+m)
}
func (m IntString) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(m))), nil
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
