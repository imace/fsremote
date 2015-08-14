package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type SogouPic struct {
	Items []struct {
		Thumburl    string `json:"thumbUrl"`
		PicURL      string `json:"pic_url"`
		Width       string `json:"width"`
		Height      string `json:"height"`
		Imgcolor    string `json:"imgcolor"`
		ThumbWidth  string `json:"thumb_width"`
		ThumbHeight string `json:"thumb_height"`
	} `json:"items"`
}

const (
	fallback        = "http://img2.funshion.com/pictures01/136/997/2/1369972.jpg"
	fallback_width  = 320
	fallback_height = 180
)

//http://pic.sogou.com/pics?query=%D1%F3%D1%F3&mood=0&picformat=0&mode=0&di=0&w=05009900&dr=1&_asf=pic.sogou.com&_ast=1438151163&dm=11&leftp=44230502&cwidth=1024&cheight=768&st=250&start=48&reqType=ajax&tn=0&reqFrom=result

func sogou_pic(hint string, width_hint, height_hint int) (uri string, width, height int) {
	uri, width, height = fallback, fallback_width, fallback_height
	sogou_search := fmt.Sprintf(`http://pic.sogou.com/pics?query=%v&mood=0&picformat=0&mode=0&di=0&w=05009900&dr=1&_asf=pic.sogou.com&_ast=1438151163&dm=11&leftp=44230502&cwidth=1024&cheight=768&st=250&start=48&reqType=ajax&tn=0&reqFrom=result`, url.QueryEscape(hint))
	resp, err := http.Get(sogou_search)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return
	}
	var v SogouPic
	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil || len(v.Items) == 0 {
		return
	}
	typ := img_dim_type(width_hint, height_hint)

	for _, img := range v.Items {
		w, h := atoi(img.ThumbWidth, 0), atoi(img.ThumbHeight, 0)
		if typ == img_dim_type(w, h) {
			uri, width, height = img.Thumburl, w, h
			break
		}
	}

	return
}

// 0 :unknown, 1 : portrait, 2 : landscape, 3 : square
func img_dim_type(w, h int) int {
	if w*h == 0 {
		return 0
	}
	v := 0
	s := w * 100 / h
	switch {
	case s < 95:
		v = 1
	case s > 105:
		v = 2
	default:
		v = 3
	}
	return v
}
