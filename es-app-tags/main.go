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

func main() {
	flag.Parse()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(xiuxiu.EsAddr))
	panic_error(err)

	xiuxiu.EsAppScan(client, xiuxiu.EsAppIndice, xiuxiu.EsAppType, func(app xiuxiu.EsApp) {
		when_es_app(client, app)
	})
}
func es_app_remove_brackets() {

}
func es_app_split_name(name string) (v []string) {
	v = append(v, strings.Replace(name, "-", "", -1))
	v = append(v, strings.Split(name, "-")...)
	return
}
func es_app_clean_name(name string) (v []string) {
	v = append(v, name)
	v = append(v, es_app_split_name(name)...)
	v = append(v, es_app_remove_brackets(name)...)
	return
}
func es_cat_expand(cat string) (v []string) {
	v = append(v, cat)
	switch cat {
	case "影视":
		v = append(v, "应用")
	case "教育":
		v = append(v, "儿童")
	case "儿歌动画":
		v = append(v, "儿童")
	}
	return
}
func when_es_app(client *elastic.Client, app xiuxiu.EsApp) {
	app.Name = strings.TrimSpace(app.Name)
	names := es_app_clean_name(app.Name)
	if len(app.TagName) > 0 {
		names = append(names, es_cat_expand(app.TagName)...)
	}
	if app.CatName == "" {
		app.CatName = "游戏"
	}

	names = append(names, es_cat_expand(app.CatName)...)
	app.Tag = xiuxiu.EsUniqSlice(names)

	if xiuxiu.EsDebug {
		fmt.Println(app.Name, app.Tag)
		return
	}
	if _, err := client.Index().
		Index(xiuxiu.EsAppIndice).
		Type(xiuxiu.EsAppType).
		Id(strconv.Itoa(app.AppID)).BodyJson(&app).Do(); err != nil {
		log.Println(err)
	} else {
		fmt.Println(app.Name)
	}
}

func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}
