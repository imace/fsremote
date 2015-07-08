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
	var doc []string
	doc = append(doc, tags(m.Name)...)

	doc = append(doc, m.Tags...)
	for _, word := range doc {
		fmt.Println(word)
	}
}
func tags(s string) []string {
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return r < 128 && !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	return fields
}
