package webapp

import (
	"aleesa-webapp-go/internal/config"
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var (
	// ForwardMax максимальное количество перенаправлений между сервисами бота, по-умолчанию, 5.
	ForwardMax int64 = 5

	// ConfigFileSize максимальная длина конфиг-файла, предполагается что 64кб хватит всем.
	ConfigFileSize int64 = 65535

	// DataDir каталог с данными, в котором размещаются например, конфиг и кшт веб-клиентов, а также ключи flickr-а.
	DataDir = "data"

	// Host хост, на котором работает сервер redis.
	Host = "localhost"

	// Loglevel уровень логгирования, по-умолчанию, info.
	Loglevel = "info"

	// RedisPort порт, на которм слушает сервер redis, по-умолчанию 6379.
	RedisPort = 6379

	// NetworkTimeout таймаут сетевого соединения.
	NetworkTimeout = 10

	// Config структурка, содержащая данные конфига.
	Config *config.MyConfig

	// ctx контекст редиски.
	Ctx = context.Background()

	// Subscriber объектик сабскрайбера редиски.
	Subscriber *redis.PubSub

	// Shutdown ставится в true, если мы получили сигнал на выключение.
	Shutdown = false

	// SigChan канал, в который приходят уведомления для хэндлера сигналов от траппера сигналов.
	SigChan = make(chan os.Signal, 1)
)

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
