package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/hearts.zhang/xiuxiu"
	"github.com/olivere/elastic"
)

func init() {

}

func main() {
	flag.Parse()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(xiuxiu.EsAddr))
	panic_error(err)

	xiuxiu.EsAppScan(client, "ottpomme2", xiuxiu.EsAppType, func(app xiuxiu.EsApp) {
		when_es_app(client, app)
	})
}
func es_app_remove_brackets(x string) string {
	re := regexp.MustCompile(`[（\(].*[\)）]`)
	return re.ReplaceAllString(x, "")

}
func es_app_split_name(name string) (v []string) {
	v = append(v, strings.Replace(name, "-", "", -1))

	v = append(v, strings.FieldsFunc(name, func(r rune) bool {
		return r == '：' || r == '-' || r == ':'
	})...)
	return
}
func es_app_clean_name(name string) (v []string) {
	v = append(v, name)

	v = append(v, es_app_split_name(name)...)
	if len(v) >= 2 {
		v = append(v, es_app_remove_upper(v[1]))
	}
	bn := es_app_remove_brackets(name)

	v = append(v, bn)
	v = append(v, es_app_remove_upper(bn))
	v = append(v, es_app_tag_from_name(name)...)
	return
}

func es_app_tag_from_name(name string) (v []string) {
	switch {
	case strings.Contains(name, "幼儿"):
		v = append(v, "儿童")
		v = append(v, "儿童教育")
		v = append(v, "教育")
	case strings.Contains(name, "浏览器"):
		if !strings.Contains(name, "文件浏览器") {
			v = append(v, "浏览器")
		}
	case strings.Contains(name, "培训"):
		v = append(v, "教育")
	case strings.Contains(name, "音乐"):
		v = append(v, "音乐")
	case strings.Contains(name, "听书"):
		v = append(v, "音乐")
		v = append(v, "听书")
	case strings.Contains(name, "助手"):
		v = append(v, "工具")
	case strings.Contains(name, "儿歌"):
		v = append(v, "儿童")
		v = append(v, "儿童教育")
		v = append(v, "教育")
	}
	return
}
func es_app_remove_upper(name string) string {
	x := strings.TrimLeftFunc(name, func(r rune) bool {
		return unicode.IsUpper(r) || unicode.IsDigit(r)
	})
	x = strings.TrimRightFunc(x, func(r rune) bool {
		return unicode.IsUpper(r) || unicode.IsDigit(r) || r == '.' || r == 'I' || r == '-' || r == '－'
	})
	return x
}

func es_cat_expand(cat string) (v []string) {
	v = append(v, cat)
	switch cat {
	case "应用":
		v = append(v, "软件")
	case "影视":
		v = append(v, "应用")
		v = append(v, "软件")
		v = append(v, "视频")
		v = append(v, "点播")
	case "教育":
		v = append(v, "儿童")
		v = append(v, "儿童教育")
	case "益智启蒙":
		v = append(v, "儿童")
		v = append(v, "儿童教育")
	case "儿童读物":
		v = append(v, "儿童")
		v = append(v, "儿童教育")
	case "儿歌动画":
		v = append(v, "儿童")
		v = append(v, "儿童教育")
	case "生活":
		v = append(v, "便捷生活")
	case "休闲":
		v = append(v, "游戏")
	case "新闻":
		v = append(v, "新闻阅读")
	}
	return
}
func when_es_app(client *elastic.Client, app xiuxiu.EsApp) {
	app.Name = strings.TrimSpace(app.Name)
	names := es_app_clean_name(app.Name)
	if len(app.TagName) > 0 {
		names = append(names, es_cat_expand(app.TagName)...)
	}

	names = append(names, es_cat_expand(app.CatName)...)
	app.Tag = xiuxiu.EsUniqSlice(names)
	app.Tags = strings.Join(app.Tag, " ")

	if xiuxiu.EsDebug {
		fmt.Println(app.Name, app.Tag)
		return
	}
	if _, err := client.Index().
		Index("ottpomme2").
		Type(xiuxiu.EsAppType).
		Id(strconv.Itoa(app.AppID)).BodyJson(&app).Do(); err != nil {
		log.Println(err)
	} else {
		fmt.Println(app.AppID)
	}
}

func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}
