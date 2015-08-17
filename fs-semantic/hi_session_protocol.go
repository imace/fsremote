package main

import (
	"encoding/json"
	"log"
)

type hi_session_protocol struct {
	Domain     string                 `json:"domain"`
	Type       string                 `json:"type"`
	OriginCode string                 `json:"originCode,omitempty"`
	OriginType string                 `json:""originType,omitempty`
	Data       map[string]interface{} `json:"data"`
	text       string
}

func when_session_protocol(p hi_session_protocol) sm_semantic {
	log.Println("what")
	_, err := json.MarshalIndent(&p, "", " ")
	panic_error(err)

	var v sm_semantic
	return v
}

//WAITING
type hi_data_waiting struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Oncancel string `json:"onCancel"`
}
type hi_data_oncancel struct {
	Domain  string `json:"domain"`
	Confirm string `json:"confirm"`
	Message string `json:"message"`
}

type hi_data_errorshow struct {
	Answer string `json:"answer"`
	Text   string `json:"text"`
	Rc     string `json:"rc"`
}
type hi_data_channelswitchshow struct {
	Value    string `json:"value"`
	Channel  string `json:"channel"`
	Text     string `json:"text"`
	Operator string `json:"operator"`
}
type hi_mutiple_contacts struct {
	Answer string `json:"answer"`
	Result struct {
		OnCancel string `json:"onCancel"`
		Person   []struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			Number     string `json:"number"`
			OnSelected string `json:"onSelected"`
			Pic        string `json:"pic"`
		} `json:"person"`
	} `json:"result"`
	Text string `json:"text"`
}
type hi_multiple_numbers struct {
	Answer string `json:"answer"`
	Result struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Numbers []struct {
			Number            string `json:"number"`
			NumberAttribution string `json:"numberAttribution"`
			OnSelected        string `json:"onSelected"`
		} `json:"numbers"`
		OnCancel string `json:"onCancel"`
		Pic      string `json:"pic"`
	} `json:"result"`
}
type hi_confirm_call struct {
	Answer string `json:"answer"`
	Result struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Number   string `json:"number"`
		OnCancel string `json:"onCancel"`
		OnOK     string `json:"onOK"`
		Pic      string `json:"pic"`
	} `json:"result"`
	Text string `json:"text"`
}

type hi_input_contact struct {
	Answer string `json:"answer"`
}

type hi_call_ok struct {
	Answer string `json:"answer"`
	Result struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Number string `json:"number"`
		Pic    string `json:"pic"`
	} `json:"result"`
}

type hi_input_freetext_sms struct {
	Answer string `json:"answer"`
	Result struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Number   string `json:"number"`
		OnCancel string `json:"onCancel"`
		OnOK     string `json:"onOK"`
		Pic      string `json:"pic"`
	} `json:"result"`
	Text string `json:"text"`
}

type hi_sms_ok struct {
	Answer string `json:"answer"`
	Result struct {
		Content string `json:"content"`
		Name    string `json:"name"`
		Number  string `json:"number"`
		Pic     string `json:"pic"`
	} `json:"result"`
}

type hi_contact_show struct {
	Answer string `json:"answer"`
	Result struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Numbers []struct {
			Number string `json:"number"`
			OnCall []struct {
				CodeUnsupportedText string `json:"codeUnsupportedText"`
				HasAnswer           bool   `json:"hasAnswer"`
				History             string `json:"history"`
				ID                  string `json:"id"`
				Rc                  string `json:"rc"`
				Semantic            struct {
					Intent struct {
						Name   string `json:"name"`
						Number string `json:"number"`
					} `json:"intent"`
				} `json:"semantic"`
				Service string `json:"service"`
			} `json:"onCall"`
			OnSMS []struct {
				Code                string `json:"code"`
				CodeUnsupportedText string `json:"codeUnsupportedText"`
				HasAnswer           bool   `json:"hasAnswer"`
				History             string `json:"history"`
				ID                  string `json:"id"`
				Rc                  string `json:"rc"`
				Semantic            struct {
					Intent struct {
						CodeUnsupportedText string `json:"codeUnsupportedText"`
						Name                string `json:"name"`
						Number              string `json:"number"`
					} `json:"intent"`
				} `json:"semantic"`
				Service         string `json:"service"`
				UnsupportedText string `json:"unsupportedText"`
			} `json:"onSMS"`
			OnSelected string `json:"onSelected"`
		} `json:"numbers"`
		OnSendContact []struct {
			CodeUnsupportedText string `json:"codeUnsupportedText"`
			HasAnswer           bool   `json:"hasAnswer"`
			History             string `json:"history"`
			ID                  string `json:"id"`
			Rc                  string `json:"rc"`
			Semantic            struct {
				Intent struct {
					Contact       string   `json:"contact"`
					ContactNumber []string `json:"contactNumber"`
				} `json:"intent"`
			} `json:"semantic"`
			Service_ string `json:"service "`
		} `json:"onSendContact"`
		Pic string `json:"pic"`
	} `json:"result"`
}

type hi_reminder_show struct {
	Answer string `json:"answer"`
	Result struct {
		Content  string `json:"content"`
		OnCancel string `json:"onCancel"`
		OnOK     string `json:"onOK"`
		Repeat   string `json:"repeat"`
		Time     string `json:"time"`
	} `json:"result"`
}

type hi_reminder_ok struct {
	Answer string `json:"answer"`
	Result struct {
		Content  string `json:"content"`
		OnCancel string `json:"onCancel"`
		OnOK     string `json:"onOK"`
		Repeat   string `json:"repeat"`
		Time     string `json:"time"`
	} `json:"result"`
	Text string `json:"text"`
}

type hi_app_launch struct {
	Answer string `json:"answer"`
	Result struct {
		Code    string `json:"code"`
		Keyword string `json:"keyword"`
		URL     string `json:"url"`
	} `json:"result"`
}

type hi_app_uninstall struct {
	Answer string `json:"answer"`
	Result struct {
		AppName     string `json:"app_name"`
		ClassName   string `json:"class_name"`
		PackageName string `json:"package_name"`
	} `json:"result"`
}

type hi_music_show struct {
	Answer string `json:"answer"`
	Result struct {
		MusicData []struct {
			Album    string `json:"album"`
			Artist   string `json:"artist"`
			Duration string `json:"duration"`
			ImageURL string `json:"imageUrl"`
			Title    string `json:"title"`
			URL      string `json:"url"`
		} `json:"musicData"`
		MusicSize string `json:"musicSize"`
	} `json:"result"`
}

type hi_channel_prog_list struct {
	Answer string `json:"answer"`
	Result struct {
		ByDate []struct {
			Date     string `json:"date"`
			Programs []struct {
				Pid   string `json:"pid"`
				Time  string `json:"time"`
				Title string `json:"title"`
			} `json:"programs"`
		} `json:"byDate"`
		Channel             string `json:"channel"`
		Code                string `json:"code"`
		CodeUnsupportedText string `json:"codeUnsupportedText"`
	} `json:"result"`
	TvChannel string `json:"tvChannel"`
}

type hi_prog_search_result struct {
	Answer string `json:"answer"`
	Result struct {
		ByDate []struct {
			Date     string `json:"date"`
			Programs []struct {
				Channel string `json:"channel"`
				Pid     string `json:"pid"`
				Time    string `json:"time"`
				Title   string `json:"title"`
			} `json:"programs"`
		} `json:"byDate"`
		Code                string `json:"code"`
		CodeUnsupportedText string `json:"codeUnsupportedText"`
	} `json:"result"`
}

type hi_prog_recommend struct {
	Answer string `json:"answer"`
	Result struct {
		Broadcasts []struct {
			Name     string `json:"name"`
			Pid      string `json:"pid"`
			Programs []struct {
				Channel string `json:"channel"`
				Pid     string `json:"pid"`
				Time    string `json:"time"`
				Title   string `json:"title"`
			} `json:"programs"`
			Score string `json:"score"`
		} `json:"broadcasts"`
		Code                string `json:"code"`
		CodeUnsupportedText string `json:"codeUnsupportedText"`
		Period              string `json:"period"`
	} `json:"result"`
}

type hi_web_show struct {
	Answer string `json:"answer"`
	Result struct {
		Keyword string `json:"keyword"`
		URL     string `json:"url"`
	} `json:"result"`
}

type hi_poi_show struct {
	Answer string `json:"answer"`
	Result struct {
		Actions []struct {
			Focus      bool   `json:"focus"`
			OnSelected string `json:"onSelected"`
			Title      string `json:"title"`
		} `json:"actions"`
		Shops []struct {
			Address    string   `json:"address"`
			BranchName string   `json:"branch_name"`
			Categories []string `json:"categories"`
			City       string   `json:"city"`
			Distance   string   `json:"distance"`
			Name       string   `json:"name"`
			Rate       string   `json:"rate"`
			Regions    []string `json:"regions"`
			Telephone  string   `json:"telephone"`
			URL        string   `json:"url"`
		} `json:"shops"`
	} `json:"result"`
	Text string `json:"text"`
}

type hi_position_show struct {
	Answer string `json:"answer"`
	Result struct {
		Address    string `json:"address"`
		City       string `json:"city"`
		Latitude   string `json:"latitude"`
		Longtitude string `json:"longtitude"`
		Position   string `json:"position"`
	} `json:"result"`
	Text string `json:"text"`
}

type hi_route_show struct {
	Answer string `json:"answer"`
	Result struct {
		FromLatitude   float64 `json:"fromLatitude"`
		FromLongtitude float64 `json:"fromLongtitude"`
		FromPosition   string  `json:"fromPosition"`
		ToCity         string  `json:"toCity"`
		ToLatitude     float64 `json:"toLatitude"`
		ToLongtitude   float64 `json:"toLongtitude"`
		ToPosition     string  `json:"toPosition"`
	} `json:"result"`
	Text string `json:"text"`
}

type hi_stock_show struct {
	Answer string `json:"answer"`
	Result struct {
		ChangeAmount        string `json:"changeAmount"`
		ChangeRate          string `json:"changeRate"`
		Code                string `json:"code"`
		CurrentPrice        string `json:"currentPrice"`
		HighestPrice        string `json:"highestPrice"`
		ImageURL            string `json:"imageUrl"`
		LowestPrice         string `json:"lowestPrice"`
		Name                string `json:"name"`
		TodayOpenPrice      string `json:"todayOpenPrice"`
		TurnOver            string `json:"turnOver"`
		UpdateTime          string `json:"updateTime"`
		YesterdayClosePrice string `json:"yesterdayClosePrice"`
	} `json:"result"`
	Text string `json:"text"`
}

type hi_translation_show struct {
	Answer string `json:"answer"`
	Result string `json:"result"`
	Text   string `json:"text"`
}

type hi_setting struct {
	Operands string `json:"operands"`
	Operator string `json:"operator"`
	Text     string `json:"text"`
}

type hi_input_freetext_weibo struct {
	Answer string `json:"answer"`
	Result struct {
		Content  string `json:"content"`
		OnCancel string `json:"onCancel"`
		OnOK     string `json:"onOK"`
		Vendor   string `json:"vendor"`
	} `json:"result"`
}

type hi_confirm_weibo struct {
	Answer string `json:"answer"`
	Result struct {
		Content string `json:"content"`
	} `json:"result"`
}

type hi_weather_show struct {
	Answer string `json:"answer"`
	Result struct {
		CityCode       string `json:"cityCode"`
		CityName       string `json:"cityName"`
		ErrorCode      int    `json:"errorCode"`
		FocusDateIndex int    `json:"focusDateIndex"`
		UpdateTime     string `json:"updateTime"`
		WeatherDays    []struct {
			CarWashIndex       string `json:"carWashIndex"`
			CurrentTemperature int    `json:"currentTemperature"`
			Day                int    `json:"day"`
			DayOfWeek          int    `json:"dayOfWeek"`
			DressIndex         string `json:"dressIndex"`
			DressIndexDesc     string `json:"dressIndexDesc"`
			HighestTemperature int    `json:"highestTemperature"`
			ImageTitleOfDay    string `json:"imageTitleOfDay"`
			ImageTitleOfNight  string `json:"imageTitleOfNight"`
			LowestTemperature  int    `json:"lowestTemperature"`
			Month              int    `json:"month"`
			Pm2_5              int    `json:"pm2_5"`
			Quality            string `json:"quality"`
			Weather            string `json:"weather"`
			Wind               string `json:"wind"`
			Year               int    `json:"year"`
		} `json:"weatherDays"`
	} `json:"result"`
	Text string `json:"text"`
}

type hi_multiple_app struct {
	Answer string `json:"answer"`
	Result struct {
		Applications []struct {
			AppName     string   `json:"appName"`
			ClassName   string   `json:"className"`
			OnSelected  []string `json:"onSelected"`
			PackageName string   `json:"packageName"`
		} `json:"applications"`
	} `json:"result"`
}

type hi_multiple_show struct {
	Answer string `json:"answer"`
	Result []struct {
		Domain  string `json:"domain"`
		OnClick string `json:"onClick"`
	} `json:"result"`
}

type hi_contact_add struct {
	Answer string `json:"answer"`
	Result struct {
		Name string `json:"name"`
	} `json:"result"`
	Text string `json:"text"`
}

type hi_sms_read struct {
	Answer string `json:"answer"`
	Text   string `json:"text"`
}

type hi_alarm_show struct {
	Answer string `json:"answer"`
	Result struct {
		OnCancel string `json:"onCancel"`
		OnOK     string `json:"onOK"`
		Repeat   string `json:"repeat"`
		Time     string `json:"time"`
	} `json:"result"`
	Text string `json:"text"`
}

type hi_alarm_ok struct {
	Answer string `json:"answer"`
	Result struct {
		Repeat string `json:"repeat"`
		Time   string `json:"time"`
	} `json:"result"`
}

type hi_talk_show struct {
	Answer string `json:"answer"`
	Text   string `json:"text"`
}
