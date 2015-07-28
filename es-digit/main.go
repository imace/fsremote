package main

//depends es-nameot
import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/hearts.zhang/xiuxiu"
	"github.com/olivere/elastic"
)

var _media xiuxiu.EsMedia

func main() {
	flag.Parse()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(xiuxiu.EsAddr))
	panic_error(err)

	xiuxiu.EsMediaScan(client, xiuxiu.EsIndice, xiuxiu.EsType, func(em xiuxiu.EsMedia) {
		when_es_media(client, em)
	})
}

func when_es_media(client *elastic.Client, em xiuxiu.EsMedia) {
	_media = em
	var names []string
	for _, name := range em.NameNorm {
		names = append(names, xiuxiu.EsNormDigit(name))
	}
	names = append(names, em.NameNorm...)
	names = xiuxiu.EsUniqSlice(names)
	if xiuxiu.EsDebug {

		fmt.Println(names)

		return
	}
	em.NameNorm = names
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
