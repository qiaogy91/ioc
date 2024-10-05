package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/qiaogy91/ioc/example/apps/app01"
	"math/rand"
	"net/http"
	"time"
)

// Gin 框架使用如下函数
// @Summary 修改文章标签
// @Description  修改文章标签
// @Tags         文章管理
// @Produce  json
// @Param id path int true "ID"
// @Param name query string true "ID"
// @Param state query int false "State"
// @Param modified_by query string true "ModifiedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /restful/v1/tags/{id} [put]
func (h *Handler) ginCreatTable(ctx *gin.Context) {
	if err := h.svc.CreatTable(ctx); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(200, "ok")
}
func (h *Handler) ginCreate(ctx *gin.Context) {
	req := &app01.CreateUserReq{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ins, err := h.svc.Create(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(200, ins)
}
func (h *Handler) ginList(ctx *gin.Context) {
	req := &app01.ListUserReq{}

	ins, err := h.svc.List(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err)
		return
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(time.Duration(r.Intn(5)) * time.Second)

	ctx.JSON(200, ins)
}
