package main

import (
	"os"
	"testing"
)

func TestFace(t *testing.T) {
	t.Skip()
	x := face_trim(face_suggest("雨", 32))
	for _, s := range x.Suggests {
		t.Log(s.Name, s.Weight)
	}

}
func _TestMain(m *testing.M) {
	load_medias()
	r := m.Run()
	os.Exit(r)
}
func TestFuzzy(t *testing.T) {
	t.Skip()
	x := fuzzy_trim(fuzzy_suggest("刘d华"))
	for _, m := range x.Suggests {
		t.Log(m.Name, m.Weight)
	}
}
