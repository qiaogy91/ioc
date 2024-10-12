package restful

import (
	"github.com/emicklei/go-restful/v3"
	"log/slog"
	"net/http"
)

func (h *Handler) HealthHandler(req *restful.Request, rsp *restful.Response) {
	res := map[string]string{"status": "ok"}

	if err := rsp.WriteHeaderAndJson(http.StatusOK, res, restful.MIME_JSON); err != nil {
		h.log.Error("send json response failed", slog.Any("err", err))
	}
}
