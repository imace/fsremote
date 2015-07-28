package main

import (
	"os"
	"testing"
)

func TestFace(t *testing.T) {
	x := face_trim(face_suggest_wrap("雨"), false)
	for _, s := range x.Suggests {
		t.Log(s.Media.Name, s.Score, s.Phrase)
	}

}
func TestMain(m *testing.M) {
	load_medias()
	r := m.Run()
	os.Exit(r)
}
