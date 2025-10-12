package webapp

import (
	"aleesa-webapp-go/internal/config"
)

var (
	ForwardMax     int64 = 5
	ConfigFileSize int64 = 65535
	DataDir              = "data"
	Host                 = "localhost"
	Loglevel             = "info"
	RedisPort            = 6379
	NetworkTimeout       = 10
	Config         *config.MyConfig
)

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
