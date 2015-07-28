package xiuxiu

import (
	"encoding/json"
	"log"
	"strings"
	"unicode"

	"github.com/olivere/elastic"
)

func EsCreateIfNotExist(client *elastic.Client, indice string) error {
	return prepare_es_index(client, indice)
}
func prepare_es_index(client *elastic.Client, indice string) (err error) {
	var b bool
	if b, err = client.IndexExists(indice).Do(); b == false && err == nil {
		err = create_index(client, indice)
	}

	return
}

func EsMediaScan(client *elastic.Client, indice, mtype string, handler func(EsMedia)) error {
	cursor, err := client.Scan(indice).Type(mtype).Size(100).Do()
	if err != nil {
		return err
	}
	for {
		result, err := cursor.Next()
		if err == elastic.EOS {
			break
		}
		panic_error(err)
		for _, hit := range result.Hits.Hits {
			var em EsMedia
			if err := json.Unmarshal(*hit.Source, &em); err != nil {
				log.Println(err)
			} else {
				handler(em)
			}
		}
	}
	return nil
}
func EsAppScan(client *elastic.Client, indice, mtype string, handler func(EsApp)) error {
	cursor, err := client.Scan(indice).Type(mtype).Size(100).Do()
	if err != nil {
		return err
	}
	for {
		result, err := cursor.Next()
		if err == elastic.EOS {
			break
		}
		panic_error(err)
		for _, hit := range result.Hits.Hits {
			var em EsApp
			if err := json.Unmarshal(*hit.Source, &em); err != nil {
				log.Println(err)
			} else {
				handler(em)
			}
		}
	}
	return nil
}
func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}
func drop_index(client *elastic.Client, index string) error {
	_, err := client.DeleteIndex(index).Do()
	return err
}

func create_index(client *elastic.Client, index string) error {
	_, err := client.CreateIndex(index).Do()
	return err
}

func strip_english(dirs []string) (names []string) {
	for _, d := range dirs {
		if len(d) > 0 && d[0] > 128 {
			names = append(names, d)
		}
	}
	return names
}

// /分隔
//去掉名字中间的·
//去掉中文名字中间的空白符
//去掉名字前后的空白符
//忽略英文名字
func EmCleanName(x string) (v []string) {
	for _, n1 := range EmSplitName(x) {
		for _, n2 := range EmSplitCnEng(n1) {
			v = append(v, EmSplitNameSep(n2)...)
		}
	}
	v = strip_english(v)
	v = uniq_string(v)
	return
}

func EmSplitName(name string) (v []string) {
	for _, name := range strings.Split(name, "/") {
		if tmp := strings.TrimSpace(name); len(tmp) > 0 {
			v = append(v, tmp)
		}
	}
	return v
}

//分隔相连的中英文
func EmSplitCnEng(name string) (v []string) {
	v = append(v, name)
	var current []rune

	var eng bool
	for _, r := range []rune(name) {
		if (r < 256 && eng == false) || (r >= 256 && eng == true) {
			if len(current) > 0 {
				if current[0] != '·' {
					v = append(v, string(current))
				}
				current = nil
				eng = !eng
			}
		}
		current = append(current, r)
	}
	if len(current) > 0 {
		v = append(v, string(current))
	}
	return
}

//去掉名字中的空白符
func string_strip_space(x string) string {
	var v []rune
	for _, r := range []rune(x) {
		if !unicode.IsSpace(r) && r != '-' && r != '、' && r != '&' && r != '#' && r != '\'' {
			v = append(v, r)
		}
	}
	return string(v)
}
func strip_sep(x string) string {
	var v []rune
	for _, r := range []rune(x) {
		if !seps(r) {
			v = append(v, r)
		}
	}
	return string(v)
}
func seps(r rune) bool {
	return unicode.IsSpace(r) || r == '·' || r == '-' || r == '：' || r == ':'
}
func EmSplitNameSep(x string) (ret []string) {
	if xr := []rune(x); len(xr) > 0 {
		if xr[0] > 256 {
			x = string_strip_space(x)
		}
		ret = append(ret, x)
		ret = append(ret, strings.FieldsFunc(x, seps)...)
		ret = append(ret, strip_sep(x))
	}
	return
}

func EsUniqSlice(slice []string) []string {
	return uniq_string(slice)
}
func uniq_string(a []string) (v []string) {
	set := make(map[string]struct{})
	for _, i := range a {
		set[i] = struct{}{}
	}
	for k := range set {
		v = append(v, k)
	}
	return v
}

//按照/分隔
func EsStringSplit(x string) (v []string) {
	fields := strings.Split(x, "/")
	for _, f := range fields {
		if x := strings.TrimSpace(f); x != "" {
			v = append(v, strings.TrimSpace(f))
		}
	}
	return
}

var a2c = map[rune]rune{
	'1': '一',
	'2': '二',
	'3': '三',
	'4': '四',
	'5': '五',
	'6': '六',
	'7': '七',
	'8': '八',
	'9': '九',
	'0': '零',
}

func EsNormDigit(n string) string {
	return norm_digit(n)
}
func norm_digit(n string) string {
	var x []rune

	for _, r := range []rune(n) {
		if t, ok := a2c[r]; ok {
			x = append(x, t)
		} else {
			x = append(x, r)
		}
	}
	return string(x)
}
