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

func EmCleanName(orig string) (names []string) {
	v := em_clean_name(orig)
	for _, name := range v {
		chars := []rune(name)
		if len(chars) > 1 && chars[0] > 256 {
			names = append(names, name)
		}
	}
	names = uniq_string(names)
	return
}
func split_cn_en(s string) (v []string) {
	var current []rune
	orig := []rune(s)
	var eng bool
	for _, r := range orig {
		if r < 256 {
			if eng == false {
				if len(current) > 0 {
					v = append(v, string(current))
					current = nil
					eng = true
				}
				current = append(current, r)
			}
		} else {
			if eng == true {
				if len(current) > 0 {
					v = append(v, string(current))
					current = nil
					eng = false
				}
				current = append(current, r)
			}
		}
	}
	if len(current) > 0 {
		v = append(v, string(current))
	}
	return
}
func strip_string_english(x string) string {
	orig := []rune(x)
	var n []rune
	for _, r := range orig {
		if r > 255 {
			n = append(n, r)
		}
	}
	return string(n)
}
func strip_english(dirs []string) (names []string) {
	for _, d := range dirs {
		if nm := strip_string_english(d); len(nm) > 0 {
			names = append(names, nm)
		}
	}
	return names
}
func EmCleanDirector(orig string) (names []string) {
	s1 := em_clean_name(orig)
	s1 = strip_english(s1)
	names = uniq_string(s1)
	return
}

// /分隔
//去掉名字中间的·
//去掉中文名字中间的空白符
//去掉名字前后的空白符
//忽略英文名字
func em_clean_name(x string) (v []string) {
	names := strings.Split(x, "/")
	for _, f := range names {
		x := strip_space(f)
		for _, xi := range x {
			if len([]rune(xi)) > 1 {
				v = append(v, xi)
			}
		}
	}
	return
}

//去掉名字中的空白符
func string_strip_space(x string) string {
	var v []rune
	for _, r := range []rune(x) {
		if !unicode.IsSpace(r) && r != '-' && r != '、' && r != '&' && r != '#' {
			v = append(v, r)
		}
	}
	return string(v)
}

func seps(r rune) bool {
	return unicode.IsSpace(r) || r == '·' || r == '-' || r == '：' || r == ':'
}
func strip_space(x string) (ret []string) {
	x = strings.TrimSpace(x)
	tmp := strings.FieldsFunc(x, seps)
	for _, w := range tmp {
		if len(w) > 0 && []rune(w)[0] > 128 { //中文
			w = string_strip_space(w)
		}
		tmp2 := strings.Fields(w)
		ret = append(ret, tmp2...)
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

func EsStringSplit(x string) (v []string) {
	fields := strings.Split(x, "/")
	for _, f := range fields {
		if x := strings.TrimSpace(f); x != "" {
			v = append(v, strings.TrimSpace(f))
		}
	}
	return
}
