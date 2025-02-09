package types

// RMsg структурка, описывающая входящее сообщение из pubsub-канала redis-ки.
type RMsg struct {
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

// SMsg структурка, описывающая исходящее сообщение в pubsub-канал redis-ки.
type SMsg struct {
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

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
