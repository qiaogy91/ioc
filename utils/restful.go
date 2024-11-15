package utils

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/gorilla/schema"
	"github.com/gorilla/websocket"
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

func SendFailed(w *restful.Response, e *ApiException) {
	if err := w.WriteHeaderAndJson(e.HttpCode, e, restful.MIME_JSON); err != nil {
		l.Error("send failed http response", slog.Any("err", err))
	}
}

func SendSuccess(w *restful.Response, v any) {
	if err := w.WriteAsJson(v); err != nil {
		l.Error("send failed http response", slog.Any("err", err))
	}
}

func WebsocketSendFailed(conn *websocket.Conn, e *ApiException) {
	if err := conn.WriteJSON(e); err != nil {
		l.Error("send failed ws response", slog.Any("err", err))
	}
}

func WebsocketSendSuccess(conn *websocket.Conn, v any) {
	if err := conn.WriteJSON(v); err != nil {
		l.Error("send success ws response", slog.Any("err", err))
	}
}
