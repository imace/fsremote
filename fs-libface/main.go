// fs-import project main.go
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/hearts.zhang/fsremote"
)

var input = flag.String("input", "e:/medias.json", "media json file")

type LibfaceLine struct {
	Frequency int
	Phrase    string
	Snippet   string
}

func main() {
	flag.Parse()
	file, err := os.Open(*input)
	panic_error(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		to_media_doc(line)
	}

}
func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}

func to_media_doc(line string) {
	var m fsremote.FunMedia
	panic_error(json.Unmarshal([]byte(line), &m))
	append_phrase(tags(m.Name), int(m.Weight*2), m.MediaId)

	append_phrase(m.Tags, int(m.Weight/5), m.MediaId)
}
func append_phrase(tags []string, weight, mediaid int) {
	for _, tag := range tags {
		append_imp(tag_suffix(tag), weight, mediaid)
		//		append_imp(tag_char(tag), weight, mediaid)
	}
}
func append_imp(tags []string, weight, mediaid int) {
	for _, tag := range tags {
		fmt.Printf("%v\t%v\t%v\n", weight, tag, mediaid)
	}
}
func tag_suffix(tag string) (v []string) {
	v = append(v, tag)
	orig := []rune(tag)
	for i := 1; i < len(orig)-1; i++ {
		suffix := orig[i:]
		v = append(v, string(suffix))
	}
	return
}
func tag_char(tag string) (v []string) {
	orig := []rune(tag)
	for _, r := range orig {
		v = append(v, string(r))
	}
	return
}
func tags(s string) []string {
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return r < 128 && !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	return fields
}
