package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/olivere/elastic"
)

type media struct {
	MediaId     int      `json:"mediaid"`
	Name        string   `json:"name"`
	NameEn      string   `json:"nameen,omitempty"`  //name_en
	NameOt      string   `json:"nameot,omitempty"`  //name_ot
	Language    string   `json:"lang,omitempty"`    //language
	MediaLength int      `json:"medialength"`       //medialength
	Country     string   `json:"country,omitempty"` //country
	Release     int64    `json:"release"`           //releasedate
	Image       string   `json:"image,omitempty"`   //imagefilepath
	Tags        []string `json:"tags,omitempty"`    //
	Weight      float64  `json:"weight"`            //
}

var input = flag.String("input", "e:/fs-movies.json", "movie json input file")

func main() {
	flag.Parse()
	client, err := elastic.NewClient(elastic.SetURL("http://testbox02.chinacloudapp.cn:9200"), elastic.SetSniff(false))
	drop_index(client, "fsmedia")
	create_index(client, "fsmedia")
	file, err := os.Open(*input)
	panic_error(err)
	defer file.Close()
	item_chan := make(chan string)
	const pcount = 8
	for i := 0; i < pcount; i++ {
		go processor(item_chan)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()

		item_chan <- txt
	}
	for i := 0; i < pcount; i++ {
		item_chan <- ""
	}
}

func processor(ic <-chan string) {
	es, err := elastic.NewClient(elastic.SetURL("http://testbox02.chinacloudapp.cn:9200"), elastic.SetSniff(false))
	if err != nil {
		log.Println(err)
	}
	for item := range ic {
		if item == "" {
			break
		}
		var m media
		if err = json.Unmarshal([]byte(item), &m); err == nil {
			es_push(es, "fsmedia", "media", m)
		}
	}
	if es != nil {
		es.Flush().Do()
	}
}

func es_push(es *elastic.Client, idx, typ string, m media) {
	if es == nil {
		return
	}

	r, err := es.Index().Index(idx).Type(typ).Id(strconv.Itoa(m.MediaId)).BodyJson(m).Do()
	log.Println(r, err)
}

func drop_index(client *elastic.Client, index string) error {
	_, err := client.DeleteIndex(index).Do()
	return err
}

func create_index(client *elastic.Client, index string) error {
	_, err := client.CreateIndex(index).Do()
	return err
}

func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}
