package main

type sm_intent struct {
	Service  string `json:"service"`
	Code     string `json:"code"`
	ImageUri string `json:"imageuri,omitempty"`
	Title    string `json:"title,omitempty"`
	Content  string `json:"content,omitempty"`
	Hint     int    `json:"hint"`
}

type sm_movie struct {
	sm_intent
}

type sm_song struct {
	sm_intent
}

type sm_m3u struct {
	sm_intent
}

type sm_app struct {
	sm_intent
}

type sm_movie_summary struct {
	sm_intent
}

type sm_app_summary struct {
	sm_intent
}

type sm_song_summary struct {
	sm_intent
}

type sm_android_intent struct {
	sm_intent
}

type sm_android_deeplink struct {
	sm_intent
}

type sm_android_setting struct {
	sm_intent
}

type sm_semantic struct {
	Movies      int
	Musics      int
	Apps        int
	TotalMovies int
	TotalMusics int
	TotalApps   int
	Intents     []interface{}
}
