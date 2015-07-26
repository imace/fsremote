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
func remove_char(tags []string) (v []string) {
	for _, tag := range tags {
		rt := []rune(tag)
		if len(rt) > 1 {
			v = append(v, tag)
		}
	}
	return
}

/*
U+4E00 - U+62FF
U+6300 - U+77FF
U+7800 - U+8CFF
U+8D00 - U+9FCC
2) 6582 characters from the CJKUI Ext A block.

Code points U+3400 to U+4DB5. Unicode 3.0 (1999).

3) 42711 characters from the CJKUI Ext B block.

Code points U+20000 to U+2A6D6. Unicode 3.1 (2001).

U+20000 - U+215FF
U+21600 - U+230FF
U+23100 - U+245FF
U+24600 - U+260FF
U+26100 - U+275FF
U+27600 - U+290FF
U+29100 - U+2A6DF
3) 4149 characters from the CJKUI Ext C block.

Code points U+2A700 to U+2B734. Unicode 5.2 (2009).

4) 222 characters from the CJKUI Ext D block.

Code points U+2B740 to U+2B81D. Unicode 6.0 (2010).

5) CJKUI Ext E block.
*/
func remove_n_chinese(tags []string) (v []string) {
	for _, tag := range tags {
		if len(tag) > 0 && []rune(tag)[0] > 0x4000 {
			v = append(v, tag)
		}
	}
	return
}
func when_es_media(client *elastic.Client, em fsremote.EsMedia) {
	x := strings.Fields(em.Tags)
	x = append(x, em.Actors...)
	x = append(x, em.Roles...)
	x = append(x, em.NameNorm...)
	x = remove_char(x)
	x = uniq_string(x)
	x = remove_n_chinese(x)
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
