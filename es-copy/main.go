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
	source, target, indice, mtype string
)

func init() {
	flag.StringVar(&source, "source", "http://172.16.13.230:9200", "or http://testbox02.chinacloudapp.cn:9200")
	flag.StringVar(&target, "target", "http://[fe80::fabc:12ff:fea2:64a6]:9200", "target indice")
	flag.StringVar(&indice, "indice", "fsmedia2", "or fsmedia")
	flag.StringVar(&mtype, "mtype", "media", "target type")
}
func prepare_es_index(client *elastic.Client, indice string) (err error) {
	var b bool
	if b, err = client.IndexExists(indice).Do(); b == false && err == nil {
		err = create_index(client, indice)
	}

	return
}
func main() {
	flag.Parse()
	src, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(source))
	panic_error(err)
	tgt, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(target))
	panic_error(err)

	err = prepare_es_index(tgt, indice)
	panic_error(err)

	cursor, err := src.Scan(indice).Type(mtype).Size(50).Do()
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

			when_es_media(tgt, em, indice)
		}
	}
}

//	_, err = client.Index().Index(es_index).Type("equip").Id(strconv.Itoa(int(e.EquipId))).BodyJson(&e).Do()
func when_es_media(client *elastic.Client, em fsremote.EsMedia, indice string) {

	if _, err := client.Index().Index(indice).Type(mtype).Id(strconv.Itoa(em.MediaID)).BodyJson(&em).Do(); err != nil {
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
	fmt.Println(em.Name, f2s(em.Weight), f2s(em.Weight2), em.MediaLength, em.Day, em.Week, em.Seven, em.Month, em.Play, em.Release)
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
