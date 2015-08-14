package xiuxiu

import (
	"strconv"
	"time"
)

type IntString int
type EsMedia struct {
	Area           string    `json:"area,omitempty"`
	Actor          string    `json:"actor"`
	Aword          string    `json:"aword"`
	Brief          string    `json:"brief,omitempty"`
	Category       string    `json:"category"`
	ChannelEn      string    `json:"channelEn"`
	Channel        string    `json:"channel"`
	Country        string    `json:"country"`
	CoverPicID     int       `json:"coverPicId"`
	PlayCountDay   int       `json:"day"`
	Director       string    `json:"director"`
	DisplayType    string    `json:"displayType"`
	Duration       string    `json:"duration,omitempty"`
	FirstCharCn    string    `json:"firstCharCn"`
	Image          string    `json:"image"`
	Lang           string    `json:"lang"`
	MediaID        int       `json:"mediaId"`
	MediaLength    IntString `json:"mediaLength"`
	Isend          IntString `json:"isend"`
	PlayCountMonth int       `json:"month"`
	Name           string    `json:"name"`
	NameEn         string    `json:"nameEn"`
	NameOt         string    `json:"nameOt"`
	Pinyin         string    `json:"pinyin,omitempty"`
	PinyinCn       string    `json:"pinyinCn,omitempty"`
	PlayCount      int       `json:"play"`
	Release        IntString `json:"release"`
	Role           string    `json:"role"`
	PlayCountSeven int       `json:"seven"`
	Tags           string    `json:"tags"`
	PlayCountWeek  int       `json:"week"`
	Weight         float64   `json:"weight"`
	Score          float64   `json:"score"`
	Actors         []string  `json:"actors,omitempty"`
	Roles          []string  `json:"roles,omitempty"`
	NameNorm       []string  `json:"namen,omitempty"`
	Directors      []string  `json:"directoren,omitempty"`
	Plots          string    `json:"plots,omitempty"`
	Digests        string    `json:"digest,omitempty"`
	Portrait       string    `json:"portrait"`
	Landscape      string    `json:"landscape"`
	ReleaseDay     string    `json:"releasedate,omitempty"`
}

func (m EsMedia) ReleaseDate() string {
	const layout = "2006-01-02"
	return time.Unix(int64(m.Release), 0).Format(layout)
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

type EsApp struct {
	AppID     int      `json:"appId"`
	CatName   string   `json:"catName"`
	DownCount int      `json:"downCount"`
	Name      string   `json:"name"`
	Pinyin    string   `json:"pinyin"`
	PkgName   string   `json:"pkgName"`
	Source    int      `json:"source"`
	TagName   string   `json:"tagName"`
	VerCode   int      `json:"verCode"`
	VerName   string   `json:"verName"`
	Tag       []string `json:"tag,omitempty"`
	Tags      string   `json:"tags,omitempty"`
	Weight    int      `json:"weight"`
}
