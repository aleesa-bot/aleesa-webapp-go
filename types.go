package main

// Конфиг
type myConfig struct {
	Server         string `json:"server,omitempty"`
	Port           int    `json:"port,omitempty"`
	Timeout        int    `json:"timeout,omitempty"`
	Loglevel       string `json:"loglevel,omitempty"`
	Log            string `json:"log,omitempty"`
	Channel        string `json:"channel,omitempty"`
	DataDir        string `json:"datadir,omitempty"`
	Csign          string `json:"csign,omitempty"`
	ForwardsMax    int64  `json:"forwards_max,omitempty"`
	OpenWeatherMap struct {
		Enabled bool   `json:"enabled,omitempty"`
		Country bool   `json:"country,omitempty"`
		Appid   string `json:"appid,omitempty"`
	} `json:"openweathermap,omitempty"`
	Flickr struct {
		Enabled bool   `json:"enabled,omitempty"`
		Key     string `json:"key,omitempty"`
		Secret  string `json:"secret,omitempty"`
	} `json:"flickr,omitempty"`
}

// Входящее сообщение из pubsub-канала redis-ки
type rMsg struct {
	From     string `json:"from,omitempty"`
	Chatid   string `json:"chatid,omitempty"`
	Userid   string `json:"userid,omitempty"`
	Threadid string `json:"threadid,omitempty"`
	Message  string `json:"message,omitempty"`
	Plugin   string `json:"plugin,omitempty"`
	Mode     string `json:"mode,omitempty"`
	Misc     struct {
		Answer      int64  `json:"answer,omitempty"`
		Botnick     string `json:"bot_nick,omitempty"`
		Csign       string `json:"csign,omitempty"`
		Fwdcnt      int64  `json:"fwd_cnt,omitempty"`
		GoodMorning int64  `json:"good_morning,omitempty"`
		Msgformat   int64  `json:"msg_format,omitempty"`
		Username    string `json:"username,omitempty"`
	} `json:"Misc"`
}

// Исходящее сообщение в pubsub-канал redis-ки
type sMsg struct {
	From     string `json:"from"`
	Chatid   string `json:"chatid"`
	Userid   string `json:"userid"`
	Threadid string `json:"threadid"`
	Message  string `json:"message"`
	Plugin   string `json:"plugin"`
	Mode     string `json:"mode"`
	Misc     struct {
		Answer      int64  `json:"answer"`
		Botnick     string `json:"bot_nick"`
		Csign       string `json:"csign"`
		Fwdcnt      int64  `json:"fwd_cnt"`
		GoodMorning int64  `json:"good_morning"`
		Msgformat   int64  `json:"msg_format"`
		Username    string `json:"username"`
	} `json:"misc"`
}

// Элемент массива, возвращемый Тhe cat api
type theCatAPIinnerStruct struct {
	Id     string `json:"id,omitempty"`
	Url    string `json:"url"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// Массив, возвращаемый на запрос урла картинки из The Cat API
type theCatAPI []theCatAPIinnerStruct

// Структура, возвращаемая на запрос в randomfox.ca API
type randomFox struct {
	Image string `json:"image,omitempty"`
	Link  string `json:"link"`
}

type obutts struct {
	Id      int    `json:"id,omitempty"`
	Author  string `json:"author,omitempty"`
	Rank    int    `json:"rank,omitempty"`
	Model   string `json:"model,omitempty"`
	Preview string `json:"preview"`
}

type oboobs struct {
	Id      int    `json:"id,omitempty"`
	Author  string `json:"author,omitempty"`
	Rank    int    `json:"rank,omitempty"`
	Model   string `json:"model,omitempty"`
	Preview string `json:"preview"`
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
