package restful

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/qiaogy91/ioc/example/apps/app01"
	"math/rand"
	"time"
)

// Restful 框架使用如下函数
func (h *Handler) restfulCreateTable(request *restful.Request, response *restful.Response) {

	if err := h.svc.CreatTable(request.Request.Context()); err != nil {
		response.WriteAsJson(err)
		return
	}
	response.WriteAsJson("ok")
}
func (h *Handler) restfulCreate(request *restful.Request, response *restful.Response) {
	req := &app01.CreateUserReq{}
	if err := request.ReadEntity(req); err != nil {
		response.WriteAsJson(err)
		return
	}
	ins, err := h.svc.Create(request.Request.Context(), req)
	if err != nil {
		response.WriteAsJson(err)
		return
	}
	response.WriteAsJson(ins)
}
func (h *Handler) restfulList(request *restful.Request, response *restful.Response) {
	req := &app01.ListUserReq{}
	ins, err := h.svc.List(request.Request.Context(), req)
	if err != nil {
		response.WriteAsJson(err)
		return
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(time.Duration(r.Intn(5)) * time.Second)
	response.WriteAsJson(ins)
}
