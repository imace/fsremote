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
			} else {
				when_es_media(client, em)
			}
		}
	}
}

//	_, err = client.Index().Index(es_index).Type("equip").Id(strconv.Itoa(int(e.EquipId))).BodyJson(&e).Do()
func when_es_media(client *elastic.Client, em fsremote.EsMedia) {
	x := strings.Fields(em.Tags)
	x = append(x, em.Actors...)
	x = append(x, em.Roles...)
	x = append(x, em.NameNorm...)
	x = uniq_string(x)
	em.Tags = strings.Join(x, " ")
	if _, err := client.Index().Index(tindice).Type(mtype).Id(strconv.Itoa(em.MediaID)).BodyJson(&em).Do(); err != nil {
		log.Println(err)
	} else {
		fmt.Println(em.MediaID)
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
