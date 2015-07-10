// fs-media project fs-media.go
package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"database/sql"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hearts.zhang/fsremote"
)

var input = flag.String("input", "e:/fs-movie-.csv", "input csv filepath")

const fs_media = "select mediaid, name_cn, name_en, name_ot, name_sn, country, medialength, releasedate, coverpicid, adword, tag4editor, tag4special, pinyin_cn, tag4spide, language from fs_media"
const fs_playnum = "select mediaid,playnum, daynum,seven_daysnum,weeknum,monthnum,modifydate from fs_media_playnum"

var _playnums = make(map[int]fsremote.FunTomato)

func play_nums() {
	db, err := sql.Open("mysql", "dbs:R4XBfuptAH@tcp(192.168.8.121:3306)/corsair_0")
	panic_error(err)
	defer db.Close()
	rows, err := db.Query(fs_playnum)
	panic_error(err)
	for rows.Next() {
		to := fsremote.FunTomato{}

		var modifydate []byte
		if err = rows.Scan(&to.MediaId, &to.PlayNum, &to.DayNum, &to.Day7Num, &to.WeekNum, &to.MonthNum, &modifydate); err == nil {
			to.Date = time_parse(string(modifydate)).Unix()
			_playnums[to.MediaId] = to
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
	}

}
func main() {
	play_nums()
	db, err := sql.Open("mysql", "dbs:R4XBfuptAH@tcp(192.168.8.121:3306)/corsair_0")
	panic_error(err)
	defer db.Close()
	rows, err := db.Query(fs_media)
	panic_error(err)
	for rows.Next() {
		var media fsremote.FunMedia
		var name_cn, name_en, name_ot, name_sn, country, medialength, releasedate,
			coverpicid, adword,
			tag4editor, tag4special, pinyin_cn, tag4spide, language []byte
		if err = rows.Scan(&media.MediaId, &name_cn,
			&name_en, &name_ot, &name_sn, &country,
			&medialength, &releasedate, &coverpicid,
			&adword, &tag4editor, &tag4special,
			&pinyin_cn, &tag4spide, &language); err == nil {
			media.Name = clean(string(name_cn))
			media.NameEn = clean(string(name_en))
			media.NameOt = clean(string(name_ot))
			media.NameSn = clean(string(name_sn))
			media.Language = clean(string(language))
			media.MediaLength = clean_medialength(string(medialength))
			media.Country = clean(string(country))
			media.Release = atoi(string(releasedate))
			media.CoverId = int(atoi(string(coverpicid)))
			media.Weight = calc_weight(media.Release, media.MediaId)
			media.Tags = append(media.Tags, tags_from(string(tag4spide))...)
			media.Tags = append(media.Tags, tags_from(string(adword))...)
			media.Tags = uniq_tags(media.Tags)

			data, _ := json.Marshal(&media)

			fmt.Println(string(data))
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}
func fill_fields(line string) map[string]int {
	v := make(map[string]int)
	f := strings.Split(line, ";")
	for i, x := range f {
		v[unquote(x)] = i
	}
	return v
}

func split_line(data []byte, eof bool) (advance int, token []byte, err error) {
	var inq bool
	for i, c := range data {
		if c == '"' {
			inq = !inq
		}
		if !inq && (c == '\r' || c == '\n') {
			advance = eat_linefeed(data, i)
			token = data[:i]
			return
		}
	}
	if eof && len(data) > 0 {
		advance = len(data)
		token = data[:]
	}
	return
}
func eat_linefeed(data []byte, idx int) int {
	for idx < len(data) && (data[idx] == '\r' || data[idx] == '\n') {
		idx++
	}
	return idx
}
func unquote(txt string) string {
	return strings.Trim(txt, `"`)
}
func mysql_spliter(data []byte, eof bool) (advance int, token []byte, err error) {
	var inq bool
	for i, c := range data {
		if c == '"' {
			inq = !inq
		}
		if !inq && c == ';' {
			advance = eat_comma(data, i)
			token = data[:i]
			return
		}
	}
	if eof && len(data) > 0 {
		advance = len(data)
		token = data[:]
	}
	return
}
func eat_comma(data []byte, idx int) int {
	for idx < len(data) && data[idx] == ';' {
		idx++
	}
	return idx
}
func mysql_split(line string) (v []string) {
	reader := strings.NewReader(line)
	scaner := NewScanner(reader)
	scaner.Split(mysql_spliter)
	for scaner.Scan() {
		v = append(v, scaner.Text())
	}

	return v
}

const (
	x1 = 1400000
	//RFC3339     = "2006-01-02T15:04:05Z07:00"
	start = "2000-01-01T00:00:00Z-08:00"
)

var start_t, _ = time.Parse(time.RFC3339, start)

func calc_weight(release int64, id int) float64 {
	days := time.Since(time.Unix(release, 0)).Hours() / 24
	if days < 1.0 {
		days = 1.0
	}
	x := _playnums[id]
	return float64(x.PlayNum)/days + float64(x.DayNum) + float64(x.Day7Num)/7.0 + float64(x.WeekNum)/5.0 + float64(x.MonthNum)/30.0
}
func uniq_tags(tags []string) []string {
	x := make(map[string]interface{})
	for _, s := range tags {
		x[s] = nil
	}
	var v []string
	for k, _ := range x {
		v = append(v, k)
	}
	return v
}
func tags_from(txt string) []string {
	txt = clean_html_tag(txt)
	txt = clean(txt)
	fields := strings.FieldsFunc(txt, func(c rune) bool {
		return c == ':' || c == ' ' || c == '：' || c == ',' || c == '，' || c == ';' || c == '；' || c == '。' || c == '\r' || c == '\n' || c == '、'
	})
	var v []string
	for _, f := range fields {
		if len(f) < 10 {
			v = append(v, f)
		}
	}
	return v
}
func clean_html_tag(txt string) string {
	var v []rune
	var ignore bool
	for _, c := range txt {
		if ignore && c != '>' {
			continue
		}
		if ignore && c == '>' {
			ignore = false
			continue
		}
		if !ignore && c == '<' {
			ignore = true
			continue
		}
		if !ignore && c != '<' {
			v = append(v, c)
		}
	}
	return string(v)
}

func clean_medialength(txt string) int {
	fields := strings.FieldsFunc(txt, func(c rune) bool { return c == ':' || c == ' ' || c == '：' })
	var x string
	if len(fields) == 1 {
		x = fields[0]
	}
	if len(fields) == 2 {
		x = fields[1]
	}
	x = strings.TrimRightFunc(x, func(c rune) bool { return !(c == '0' || (c >= '1' && c <= '9')) })
	v, _ := strconv.Atoi(x)
	return v
}
func clean_nonchar(r rune) (v rune) {
	v = -1
	switch r {
	case ';', ',', '/', '\'', '"', '。', '，', '·', '《', '》', '“', '”', '（', '）', '【', '】':
		v = ' '
	default:
		v = r
	}
	return
}
func clean(txt string) string {
	v := strings.Map(clean_nonchar, txt)
	return strings.Join(strings.Fields(v), " ")
}

func atoi(i string) int64 {
	v, _ := strconv.ParseInt(i, 10, 64)
	return v
}

//http://img[1-4].funshion.com/pictures01/420/454/420454.jpg
//http://img[1-4].funshion.com/pictures/987/654/3/9876543.jpg
func image_url(id string, portrait bool) string {
	oid := id
	var path []string
	for l := len(id); l >= 3; {
		path = append(path, id[:3])
		id = id[3:]
		l = len(id)
	}
	if len(id) > 0 {
		path = append(path, id)
	}
	root := "pictures"
	if !portrait {
		root = root + "01"
	}
	idx := rand.Intn(4) + 1
	return "http://img" + strconv.Itoa(idx) + ".funshion.com/" + root + "/" + strings.Join(path, "/") + "/" + oid + ".jpg"
}

const tmlayout = "2006-01-02 15:04:05 -0700"

func time_parse(t string) time.Time {
	v, _ := time.Parse(tmlayout, t+" +0800")
	return v
}

//"mediaid";"name_cn";"name_en";"name_ot";"name_sn";"language";"medialength";"country";"tv_station";"website";"releasedate";"releaseinfo";"firstrun_date";"imagefilepath";"coverpicid";"adword";"behind";"plots";"firstchar_cn";"firstchar_en";"tag4editor";"tag4special";"displaytype";"ordering";"isplay";"isdisplay";"statdate";"createdate";"modifydate";"createuserid";"modifyuserid";"relatedmedia";"relatedmedia_editor";"relatedmedia_mobile";"pos";"ad";"webplay";"webclarity";"pcclarity";"isrank";"isclassic";"istheatre";"iscutpic";"special_cutpic";"tactics";"issue";"torrentnum";"torrentall";"torrentdesc";"fsp_status";"fsp_lang_status";"fsp_original_status";"fsp_info";"pinyin_cn";"supporttype";"ta_0";"ta_1";"ta_2";"ta_3";"ta_4";"ta_5";"ta_6";"ta_7";"ta_8";"ta_9";"copyright";"deleted";"extinfo";"program_type";"progpicdate";"eventlink";"update_pattern";"info_modifydate";"media_index";"awards";"auto_publish";"version";"hasvtags";"top_pic_timestamp";"timelen";"tag4spide";"nami_state"
