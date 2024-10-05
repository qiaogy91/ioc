package restful

import (
	"github.com/emicklei/go-restful/v3"
	"net/http"
)

func (h *Handler) HealthHandler(req *restful.Request, rsp *restful.Response) {
	res := map[string]string{"status": "ok"}

	if err := rsp.WriteHeaderAndJson(http.StatusOK, res, restful.MIME_JSON); err != nil {
		h.log.Error().Msgf("send json response failed, %s", err)
	}
}
