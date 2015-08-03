package intent

const (
	card_type_simple          = "card.fun.simple"
	card_type_app_open        = "card.fun.app.open"
	card_type_app_uninstall   = "card.fun.app.uninstall"
	card_type_video           = "card.fun.app.video"
	card_type_music           = "card.fun.app.music"
	card_type_tv_operator     = "card.fun.tv.operator"
	card_type_search_generic  = "card.fun.search.generic"
	card_type_search_movie    = "card.fun.search.movie"
	card_type_recipe          = "card.fun.recipe"
	card_type_news            = "card.fun.news"
	card_type_setting_generic = "card.fun.setting.generic"
	card_type_tv_setting      = "card.fun.tv.setting"
	card_type_website         = "card.fun.website"
	card_type_stock           = "card.fun.stock"
	card_type_weather         = "card.fun.weather"
)

type IntentCard struct {
	Title string      `json:"title,omitempty"`
	Desc  string      `json:"desc,omitempty"`
	Image string      `json:"image,omitempty"`
	Typ   string      `json:"type,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}
