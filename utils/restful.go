package utils

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/gorilla/schema"
	"github.com/qiaogy91/ioc/config/log"
	"log/slog"
)

func init() {
	l = log.Sub("utils")
	Decoder = schema.NewDecoder()
}

var (
	l       *slog.Logger
	Decoder *schema.Decoder
)

func SendFailed(w *restful.Response, e error) {
	if err := w.WriteAsJson(e); err != nil {
		l.Error("send failed response err", slog.Any("err", err))
	}
}

func SendSuccess(w *restful.Response, v any) {
	if err := w.WriteAsJson(v); err != nil {
		l.Error("send failed response err", slog.Any("err", err))
	}
}
