package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hearts.zhang/xiuxiu"
	"github.com/olivere/elastic"
)

func init() {

}

var _media xiuxiu.EsMedia

func main() {
	flag.Parse()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(xiuxiu.EsAddr))
	panic_error(err)

	xiuxiu.EsMediaScan(client, xiuxiu.EsIndice, xiuxiu.EsType, func(em xiuxiu.EsMedia) {
		when_es_media(client, em)
	})
}

func remove_china_hk(tags []string, loc string) (v []string) {
	if !strings.Contains(loc, "中国香港") {
		return tags
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
		tags = append(tags, "大陆")
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
		tags = append(tags, "日韩")
	}
	if strings.Contains(loc, "韩国") {
		tags = append(tags, "韩剧")
		tags = append(tags, "日韩")
	}

	return
}
func type2tag(typ string) (tags []string) {
	switch typ {
	case "movie":
		tags = append(tags, "电影")
		tags = append(tags, "影片")
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
	tags = xiuxiu.EsStringSplit(loc)
	if strings.Contains(loc, "中国") && !strings.Contains(loc, "香港") {
		tags = append(tags, "国产")
		tags = append(tags, "国内")
		tags = append(tags, "大陆")
	}
	if strings.Contains(loc, "中国台湾") {
		tags = append(tags, "台湾")
	}
	if strings.Contains(loc, "中国香港") {
		tags = append(tags, "港片")
		tags = append(tags, "香港")
	} else if strings.Contains(loc, "香港") {
		tags = append(tags, "港片")
	}
	if strings.Contains(loc, "英国") {
		tags = append(tags, "英剧")
	}
	if strings.Contains(loc, "日本") {
		tags = append(tags, "日韩")
		tags = append(tags, "日剧")
	}
	if strings.Contains(loc, "韩国") {
		tags = append(tags, "日韩")
		tags = append(tags, "韩剧")
	}

	return
}
func remove_char(tags []string) (v []string) {
	for _, tag := range tags {
		rt := []rune(tag)
		if len(rt) > 1 || (len(rt) == 1 && rt[0] > 256) { // exclude ascii char
			v = append(v, tag)
		} else {
			log.Println(tag)
		}
	}
	return
}

func when_es_media(client *elastic.Client, em xiuxiu.EsMedia) {
	_media = em
	tags := strings.Fields(em.Tags)

	tags = append(tags, type2tag(em.DisplayType)...)
	tags = append(tags, location2tag(em.Country)...)
	tags = append(tags, loctyp2tag(em.Country, em.DisplayType)...)
	tags = remove_char(tags)
	tags = xiuxiu.EsUniqSlice(tags)
	tags = remove_china_hk(tags, em.Country)
	if xiuxiu.EsDebug {
		fmt.Println(em.Name, tags)
		fmt.Println(em.Name, em.Tags)
		return
	}
	em.Tags = strings.Join(tags, " ")
	if _, err := client.Index().Index(xiuxiu.EsIndice).Type(xiuxiu.EsType).Id(strconv.Itoa(em.MediaID)).BodyJson(&em).Do(); err != nil {
		log.Println(err)
	} else {
		fmt.Println(em.MediaID)
	}
}

func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}
