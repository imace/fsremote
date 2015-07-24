package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hearts.zhang/fsremote"
	"github.com/olivere/elastic"
)

var (
	es      string
	tindice string
)

const (
	indice = "fsmedia2"
	mtype  = "media"
)

func init() {
	flag.StringVar(&es, "es", "http://[fe80::fabc:12ff:fea2:64a6]:9200", "or http://testbox02.chinacloudapp.cn:9200")
	flag.StringVar(&tindice, "indice", "fsmedia2", "target indice")
}
func prepare_es_index(client *elastic.Client) (err error) {
	var b bool
	if b, err = client.IndexExists(tindice).Do(); b == false && err == nil {
		err = create_index(client, tindice)
	}

	return
}

var _media fsremote.EsMedia

func main() {
	flag.Parse()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(es))
	panic_error(err)

	err = prepare_es_index(client)
	panic_error(err)

	cursor, err := client.Scan(indice).Type(mtype).Size(20).Do()
	panic_error(err)
	for {
		result, err := cursor.Next()
		if err == elastic.EOS {
			break
		}
		panic_error(err)
		for _, hit := range result.Hits.Hits {
			var em fsremote.EsMedia
			if err := json.Unmarshal(*hit.Source, &em); err != nil {
				log.Println(err)
			}

			when_es_media(client, em)
		}
	}
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
func remove_china_hk(tags []string, loc string) (v []string) {
	if !strings.Contains(loc, "中国香港") {
		return
	}
	for _, tag := range tags {
		switch tag {
		case "内地剧":
		case "大陆剧":
		case "港剧":
		case "国产":
		case "国内":
		case "大陆":
		case "中国":
		case "港片":
		case "港片儿":

		default:
			v = append(v, tag)
		}
	}
	return
}
func loctyp2tag(loc, typ string) (tags []string) {
	if typ == "movie" && strings.Contains(loc, "美国") {
		tags = append(tags, "好莱坞")

	}
	if typ != "tv" {
		return
	}

	if strings.Contains(loc, "中国") && !strings.Contains(loc, "香港") {
		tags = append(tags, "内地剧")
		tags = append(tags, "大陆剧")
	}

	if strings.Contains(loc, "美国") {
		tags = append(tags, "美剧")
	}
	if strings.Contains(loc, "香港") {
		tags = append(tags, "港剧")
	}
	if strings.Contains(loc, "英国") {
		tags = append(tags, "英剧")
	}
	if strings.Contains(loc, "日本") {
		tags = append(tags, "日剧")
	}
	if strings.Contains(loc, "韩国") {
		tags = append(tags, "韩剧")
	}
	if strings.Contains(loc, "韩国") && strings.Contains(loc, "日本") {
		tags = append(tags, "日韩")
	}
	return
}
func type2tag(typ string) (tags []string) {
	switch typ {
	case "movie":
		tags = append(tags, "电影")
		tags = append(tags, "影片")
		tags = append(tags, "影片儿")
		tags = append(tags, "片儿")
	case "tv":
		tags = append(tags, "电视")
		tags = append(tags, "电视剧")
		tags = append(tags, "连续剧")
		tags = append(tags, "系列剧")
	case "cartoon":
		tags = append(tags, "动漫")
		tags = append(tags, "动画")
		tags = append(tags, "卡通")
	case "viriety":
		tags = append(tags, "综艺")
		tags = append(tags, "娱乐")
	case "vfilm":
		tags = append(tags, "微电影")
		tags = append(tags, "微视频")
		tags = append(tags, "微视")
		tags = append(tags, "小视频")
		tags = append(tags, "小电影")
	}
	return
}

func location2tag(loc string) (tags []string) {
	tags = em_split_string(loc)
	if strings.Contains(loc, "中国") {
		tags = append(tags, "国产")
		tags = append(tags, "国内")
		tags = append(tags, "大陆")
	}

	if strings.Contains(loc, "香港") {
		tags = append(tags, "港片")
		tags = append(tags, "港片儿")
	}
	if strings.Contains(loc, "英国") {
		tags = append(tags, "英剧")
	}
	if strings.Contains(loc, "日本") {
		tags = append(tags, "日剧")
	}
	if strings.Contains(loc, "韩国") {
		tags = append(tags, "韩剧")
	}

	return
}

//	_, err = client.Index().Index(es_index).Type("equip").Id(strconv.Itoa(int(e.EquipId))).BodyJson(&e).Do()
func when_es_media(client *elastic.Client, em fsremote.EsMedia) {
	_media = em
	tags := strings.Fields(em.Tags)
	tags = remove_china_hk(tags, em.Country)
	tags = append(tags, type2tag(em.DisplayType)...)
	tags = append(tags, location2tag(em.Country)...)
	tags = append(tags, loctyp2tag(em.Country, em.DisplayType)...)
	tags = uniq_string(tags)
	em.Tags = strings.Join(tags, " ")
	if _, err := client.Index().Index(tindice).Type(mtype).Id(strconv.Itoa(em.MediaID)).BodyJson(&em).Do(); err != nil {
		log.Println(err)
	} else {
		fmt.Println(em.MediaID)
	}
}
func em_split_string(x string) (v []string) {
	fields := strings.Split(x, "/")
	for _, f := range fields {
		if x := strings.TrimSpace(f); x != "" {
			v = append(v, strings.TrimSpace(f))
		}
	}
	return
}
func print_es_media(em fsremote.EsMedia) {
	fmt.Println(em.Name, f2s(em.Weight), em.MediaLength, em.Day, em.Week, em.Seven, em.Month, em.Play, em.Release)
}
func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}
func f2s(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 3, 64)
}
func drop_index(client *elastic.Client, index string) error {
	_, err := client.DeleteIndex(index).Do()
	return err
}

func create_index(client *elastic.Client, index string) error {
	_, err := client.CreateIndex(index).Do()
	return err
}
