package main

// myConfig описывает структуру Конфига программы.
type myConfig struct {
	Loglevel string `json:"loglevel,omitempty"`
	DataDir  string `json:"datadir,omitempty"`
}

// monkeyUsers описывает структуру json-массива ответа www.monkeyusers.com.
type monkeyUsers []struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
