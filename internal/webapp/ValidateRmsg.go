package webapp

import (
	"aleesa-webapp-go/internal/config"
	"fmt"
)

// ValidateRmsg валидирует входящее сообщение.
func ValidateRmsg(cfg *config.MyConfig, j RMsg, msg string) (RMsg, error) {
	if exist := j.From; exist == "" {
		emsg := fmt.Errorf("incorrect msg from redis, no from field: %s", msg)

		return j, emsg
	}

	if exist := j.Chatid; exist == "" {
		emsg := fmt.Errorf("incorrect msg from redis, no chatid field: %s", msg)

		return j, emsg
	}

	if exist := j.Userid; exist == "" {
		emsg := fmt.Errorf("incorrect msg from redis, no userid field: %s", msg)

		return j, emsg
	}

	// j.Threadid может быть пустым, значит либо нам его не дали, либо дали пустым. Это нормально.

	if exist := j.Message; exist == "" {
		emsg := fmt.Errorf("incorrect msg from redis, no message field: %s", msg)

		return j, emsg
	}

	if exist := j.Plugin; exist == "" {
		emsg := fmt.Errorf("incorrect msg from redis, no plugin field: %s", msg)

		return j, emsg
	}

	if exist := j.Mode; exist == "" {
		emsg := fmt.Errorf("incorrect msg from redis, no mode field: %s", msg)

		return j, emsg
	}

	// j.Misc.Answer может и не быть, тогда ответа на такое сообщение не будет
	// j.Misc.Botnick тоже можно не передавать, тогда будет записана пустая строка
	// j.Misc.Csign если нам его не передали, возьмём значение из конфига
	if exist := j.Misc.Csign; exist == "" {
		j.Misc.Csign = cfg.Csign
	}

	// j.Misc.Fwdcnt если нам его не передали, то будет 0
	if exist := j.Misc.Fwdcnt; exist == 0 {
		j.Misc.Fwdcnt = 1
	}

	// j.Misc.GoodMorning может быть быть 1 или 0, по-умолчанию 0
	// j.Misc.Msgformat может быть быть 1 или 0, по-умолчанию 0
	// j.Misc.Username можно не передавать, тогда будет пустая строка
	return j, nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
