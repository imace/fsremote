package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	front string
	dic   string
	word  string
)

func init() {
	flag.StringVar(&front, "es", "172.16.13.230:80", "php wraped es address")
	flag.StringVar(&dic, "dic", "", "volcabulary")
	flag.StringVar(&word, "word", "刘德华", "search word")
}
func search(word string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("error", word, r)
		}
	}()
	uri := es_media_url(word)
	fmt.Println(uri)
	resp, err := http.Get(uri)
	panic_error(err)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("http-error ", word, resp.Status)
		return
	}
	var medias Medias
	err = json.NewDecoder(resp.Body).Decode(&medias)
	panic_error(err)

	fmt.Printf("%v total %v match\n", word, medias.Num)
	for _, m := range medias.Data {
		fmt.Printf("%v %.2f %.2f %v %v\n", m.Name, m.Weight, m.Xscore, tm(m.Release), m.Tags)
	}
}
func main() {
	flag.Parse()
	if dic != "" {
		scan_and_search(dic)
	} else {
		search(word)
	}
}

func tm(t int64) string {
	tx := time.Unix(t, 0)
	const layout = "2006-01-02"
	return tx.Format(layout)
}
func scan_and_search(dic string) {
	file, err := os.Open(dic)
	panic_error(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) > 0 {
			search(fields[0])
		}
		fmt.Println()
	}
}

func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}
func es_media_url(term string) string {
	params := url.Values{"tags": []string{term}, "to": []string{"10"}}
	//http://es.fun.tv/media?tags=蝙蝠侠&from=0&to=100
	return "http://" + front + "/media?" + params.Encode()
}

type IntString int

func (m IntString) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(m))), nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *IntString) UnmarshalJSON(data []byte) error {
	t, err := strconv.Unquote(string(data))
	if err != nil {
		t = string(data)
	}
	i, err := strconv.Atoi(t)
	*m = IntString(i)
	return err
}

type Medias struct {
	Data []struct {
		Xscore      float64   `json:"_score"`
		Actor       string    `json:"actor"`
		Actors      string    `json:"actors"`
		Area        string    `json:"area"`
		Aword       string    `json:"aword"`
		Category    string    `json:"category"`
		Channel     string    `json:"channel"`
		ChannelEn   string    `json:"channelEn"`
		Country     string    `json:"country"`
		CoverPicID  int       `json:"coverPicId"`
		Day         int       `json:"day"`
		Director    string    `json:"director"`
		DisplayType string    `json:"displayType"`
		FirstCharCn string    `json:"firstCharCn"`
		Image       string    `json:"image"`
		Isend       IntString `json:"isend"`
		Landscape   string    `json:"landscape"`
		Lang        string    `json:"lang"`
		MediaID     int       `json:"mediaId"`
		MediaLength int       `json:"mediaLength"`
		Month       int       `json:"month"`
		Name        string    `json:"name"`
		NameEn      string    `json:"nameEn"`
		NameOt      string    `json:"nameOt"`
		Namen       string    `json:"namen"`
		Pinyin      string    `json:"pinyin"`
		PinyinCn    string    `json:"pinyinCn"`
		Play        int       `json:"play"`
		Plots       string    `json:"plots"`
		Portrait    string    `json:"portrait"`
		Release     int64     `json:"release"`
		Role        string    `json:"role"`
		Roles       string    `json:"roles"`
		Score       float64   `json:"score"`
		Seven       int       `json:"seven"`
		Tags        string    `json:"tags"`
		Week        int       `json:"week"`
		Weight      float64   `json:"weight"`
	} `json:"data"`
	Num int `json:"num"`
}
