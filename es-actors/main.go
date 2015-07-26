package main

//depends nil
import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/hearts.zhang/xiuxiu"
	"github.com/olivere/elastic"
)

func init() {

}

func main() {
	flag.Parse()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(xiuxiu.EsAddr))
	panic_error(err)

	err = xiuxiu.EsCreateIfNotExist(client, xiuxiu.EsIndice)
	panic_error(err)
	xiuxiu.EsMediaScan(client, xiuxiu.EsIndice, xiuxiu.EsType, func(em xiuxiu.EsMedia) {
		when_es_media(client, em)
	})
}

func when_es_media(client *elastic.Client, em xiuxiu.EsMedia) {
	em.Actors = xiuxiu.EmCleanName(em.Actor)
	em.Roles = xiuxiu.EmCleanName(em.Role)

	if xiuxiu.EsDebug {
		if len(em.Actors) > 0 {
			fmt.Println(em.Actors, em.Actor)
		}
		if len(em.Roles) > 0 {
			fmt.Println(em.Roles, em.Role)
		}
		return
	}
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
