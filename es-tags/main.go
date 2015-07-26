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

func main() {
	flag.Parse()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(xiuxiu.EsAddr))
	panic_error(err)

	xiuxiu.EsMediaScan(client, xiuxiu.EsIndice, xiuxiu.EsType, func(em xiuxiu.EsMedia) {
		when_es_media(client, em)
	})
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

func remove_n_chinese(tags []string) (v []string) {
	for _, tag := range tags {
		if len(tag) > 0 && []rune(tag)[0] > 0x4000 {
			v = append(v, tag)
		}
	}
	return
}
func when_es_media(client *elastic.Client, em xiuxiu.EsMedia) {
	x := strings.Fields(em.Tags)
	x = append(x, em.Actors...)
	x = append(x, em.Roles...)
	x = append(x, em.NameNorm...)
	x = remove_char(x)
	x = remove_n_chinese(x)
	x = xiuxiu.EsUniqSlice(x)

	em.Tags2 = strings.Join(x, " ")
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

func print_es_media(em xiuxiu.EsMedia) {
	fmt.Println(em.Name, f2s(em.Weight), em.MediaLength, em.Day, em.Week, em.Seven, em.Month, em.Play, em.Release)
}

func f2s(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 3, 64)
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
