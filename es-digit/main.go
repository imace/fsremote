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

var a2c = map[rune]rune{
	'1': '一',
	'2': '二',
	'3': '三',
	'4': '四',
	'5': '五',
	'6': '六',
	'7': '七',
	'8': '八',
	'9': '九',
	'0': '零',
}

func norm_digit(n string) string {
	var x []rune

	for _, r := range []rune(n) {
		if t, ok := a2c[r]; ok {
			x = append(x, t)
		} else {
			x = append(x, r)
		}
	}
	return string(x)
}

var digs = map[rune]struct{}{
	'一': struct{}{},
	'二': struct{}{},
	'三': struct{}{},
	'四': struct{}{},
	'五': struct{}{},
	'六': struct{}{},
	'七': struct{}{},
	'八': struct{}{},
	'九': struct{}{},
}

func when_es_media(client *elastic.Client, em xiuxiu.EsMedia) {
	_media = em
	var names []string
	for _, name := range em.NameNorm {
		names = append(names, norm_digit(name))
	}
	names = append(names, em.NameNorm...)
	names = xiuxiu.EsUniqSlice(names)
	if xiuxiu.EsDebug {
		if len(names) != len(em.NameNorm) {
			fmt.Println(names)
		}
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
