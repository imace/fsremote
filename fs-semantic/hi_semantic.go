package main

import "fmt"

type hi_understand_result struct {
	Code    string `json:"code"`
	History string `json:"history"`
	Text    string `json:"text"`
	Service string `json:"service"`
	Rc      int    `json:"rc"`
	Error   *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
	Semantic *struct {
		Intent          map[string]interface{} `json:"intent"`
		NormalHeader    string                 `json:"normalHeader,omitempty"`
		NormalHeaderTts string                 `json:"normalHeaderTts,omitempty"`
		NoDataHeader    string                 `json:"noDataHeader,omitempty"`
	} `json:"semantic,omitempty"`
	Data *struct {
		Header    string                 `json:"header"`
		HeaderTts string                 `json:"headerTts,omitempty"`
		Result    map[string]interface{} `json:"result,omitempty"`
	} `json:"data,omitempty"`
	General *struct {
		Type    string `json:"type"` //T,U,I TU,IT, ITU
		Text    string `json:"text,omitempty"`
		TextTts string `json:"textTts,omitempty"`
		ImgUrl  string `json:"imgUrl,omitempty"`
		ImgAlt  string `json:"imgAlt,omitempty"`
		Url     string `json:"url,omitempty"`
		UrlAlt  string `json:"urlAlt,omitempty"`
	} `json:"general,omitempty"`
}

func when_understand(result hi_understand_result) {
	fmt.Println(result.Code, result.Service)
}

type his_setting_tv struct {
	intent map[string]interface{}
}

type hii_setting_tv struct {
	operator string
	operands string
	value    string
}

const (
	sv_alarm        = "cn.yunzhisheng.alarm"        // 闹钟 闹钟
	sv_appmgr       = "cn.yunzhisheng.appmgr"       // 应用 应用管理
	sv_broadcast    = "cn.yunzhisheng.broadcast"    // 广播 电台广播
	sv_calculator   = "cn.yunzhisheng.calculator"   // 计算器 计算器
	sv_calendar     = "cn.yunzhisheng.calendar"     // 日历 日历相关问答
	sv_call         = "cn.yunzhisheng.call"         // 电话 打电话
	sv_calltransfer = "cn.yunzhisheng.calltransfer" // 呼叫转移 呼叫转移
	sv_contact      = "cn.yunzhisheng.contact"      // 联系人 新建，查找，发送联系人信息
	sv_cookbook     = "cn.yunzhisheng.cookbook"     // 菜谱 菜谱
	sv_ecommerce    = "cn.yunzhisheng.ecommerce"    // 购物 电商购物搜索
	sv_localsearch  = "cn.yunzhisheng.localsearch"  // 周边 查找周边，城市或者区域内的商家，或 POI 信息，以及周边的优惠或团购信息
	sv_flight       = "cn.yunzhisheng.flight"       // 航班 航班时刻信息
	sv_hotline      = "cn.yunzhisheng.hotline"      // 常用号码 常用号码
	sv_map          = "cn.yunzhisheng.map"          // 定位 定位和路线
	sv_microblog    = "cn.yunzhisheng.microblog"    // 微博 发微博，新浪微博、腾讯微博、人人网等
	sv_movie        = "cn.yunzhisheng.movie"        // 电影 影讯(电影上映信息)
	sv_music        = "cn.yunzhisheng.music"        // 音乐 音乐搜索
	sv_news         = "cn.yunzhisheng.news"         // 新闻 新闻搜索
	sv_note         = "cn.yunzhisheng.note"         // 备忘 备忘
	sv_novel        = "cn.yunzhisheng.novel"        // 小说 小说搜索
	sv_reminder     = "cn.yunzhisheng.reminder"     // 提醒 提醒
	sv_recharge     = "cn.yunzhisheng.recharge"     // 充值 手机充值操作
	sv_restaurant   = "cn.yunzhisheng.restaurant"   // 点餐 点餐操作
	sv_setting      = "cn.yunzhisheng.setting"      // 设置 手机等设备相关的设置
	sv_setting_tv   = "cn.yunzhisheng.setting.tv"   // 电视操作 电视机相关的设置和操作
	sv_sms          = "cn.yunzhisheng.sms"          // 短信 发短信，查看短信
	sv_stock        = "cn.yunzhisheng.stock"        // 股价 股价查询
	sv_traffic      = "cn.yunzhisheng.traffic"      // 路况 路况查询
	sv_train        = "cn.yunzhisheng.train"        // 火车 火车时刻信息
	sv_translation  = "cn.yunzhisheng.translation"  // 翻译 在线翻译
	sv_tv           = "cn.yunzhisheng.tv"           // 电视 电视节目信息
	sv_video        = "cn.yunzhisheng.video"        // 视频 视频搜索
	sv_weather      = "cn.yunzhisheng.weather"      // 天气 天气
	sv_websearch    = "cn.yunzhisheng.websearch"    // 搜索 Web 搜索
	sv_website      = "cn.yunzhisheng.website"      // 网站 网站导航
)

const (
	hicode_alarm_set              = "ALARM_SET"
	hicode_app_launch             = "APP_LAUNCH"
	hicode_app_uninstall          = "APP_UNINSTALL"
	hicode_app_download           = "APP_DOWNLOAD"
	hicode_app_install            = "APP_INSTALL"
	hicode_app_search             = "APP_SEARCH"
	hicode_app_exit               = "APP_EXIT"
	hicode_forecast               = "FORECAST"
	hicode_query                  = "QUERY"
	hicode_call                   = "CALL"
	hicode_redial                 = "REDIAL"
	hicode_reply                  = "REPLY"
	hicode_call_through           = "CALL_THROUGH"
	hicode_contact_add            = "CONTACT_ADD"
	hicode_contact_search         = "CONTACT_SEARCH"
	hicode_contact_send           = "CONTACT_SEND"
	hicode_search                 = "SEARCH"
	hicode_business_search        = "BUSINESS_SEARCH"
	hicode_deal_search            = "DEAL_SEARCH"
	hicode_nonbusiness_search     = "NONBUSINESS_SEARCH"
	hicode_flight_oneway          = "FLIGHT_ONEWAY"
	hicode_flight_twoway          = "FLIGHT_TWOWAY"
	hicode_flight_info            = "FLIGHT_INFO"
	hicode_search_song            = "SEARCH_SONG"
	hicode_search_artist          = "SEARCH_ARTIST"
	hicode_search_random          = "SEARCH_RANDOM"
	hicode_search_billboard       = "SEARCH_BILLBOARD "
	hicode_note_record            = "NOTE_RECORD"
	hicode_reminder_set           = "REMINDER_SET"
	hicode_reminder_remove        = "REMINDER_REMOVE"
	hicode_position               = "POSITION"
	hicode_route                  = "ROUTE"
	hicode_explicit_search        = "EXPLICIT_SEARCH "
	hicode_maybe_search           = "MAYBE_SEARCH "
	hicode_setting_exec           = "SETTING_EXEC"
	hicode_setting_exec_tv        = "SETTING_EXEC_TV"
	hicode_website_open           = "WEBSITE_OPEN"
	hicode_sms_send               = "SMS_SEND"
	hicode_sms_read               = "SMS_READ"
	hicode_stock_info             = "STOCK_INFO"
	hicode_train_oneway           = "TRAIN_ONEWAY"
	hicode_train_info             = "TRAIN_INFO"
	hicode_translation            = "TRANSLATION"
	hicode_play                   = "PLAY"
	hicode_play_syn               = "PLAY_SYN"
	hicode_send_to_stb            = "SEND_TO_STB"
	hicode_microblog_send         = "MICROBLOG_SEND"
	hicode_yp_call                = "YP_CALL"
	hicode_yp_search              = "YP_SEARCH"
	hicode_answer                 = "ANSWER"
	hicode_channel_prog_list      = "CHANNEL_PROG_LIST"
	hicode_prog_search_result     = "PROG_SEARCH_RESULT"
	hicode_prog_recommend         = "PROG_RECOMMEND"
	hicode_recharge_myself        = "RECHARGE_MYSELF"
	hicode_recharge_others        = "RECHARGE_OTHERS"
	hicode_traffic_query          = "TRAFFIC_QUERY"
	hicode_open                   = "OPEN"
	hicode_close                  = "CLOSE"
	hicode_product_search_keyword = "PRODUCT_SEARCH_KEYWORD"
	hicode_app_help               = "APP_HELP"
)

var errors = map[int]string{
	2010: "appkey 缺失或无效",
	2011: "appkey 有效期已过",
	2020: "应用请求签名缺失或错误",
	2030: "method 参数缺失或无效",
	2040: "text 参数缺失或无效",
	2041: "text 参数值超出长度限制",
	2042: "text 参数值长度小于最少限制",
	2050: "ver 参数缺失或无效",
	2060: "time 参数无效",
	3010: "IP 被列为黑名单状态",
	3011: "Appkey 访问总数超限",
	3012: "Appkey 访问频率超限",
}

/*
{
    "code": "SETTING_EXEC_TV",
    "history": "cn.yunzhisheng.setting.tv",
    "text": "看中央1",
    "semantic": {
    "intent": {
    "operator": "ACT_OPEN_CHANNEL",
    "operands": "CCTV-1",
    "value": "中央一"
    }
    },
    "service": "cn.yunzhisheng.setting.tv",
    "rc": 0
    }
*/
