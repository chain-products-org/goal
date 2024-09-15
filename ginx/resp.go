package ginx

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/gophero/goal/errorx"
)

var Resp = new(resp)

type resp struct{}

func (r *resp) Ok(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

func (r *resp) OkJson(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, data)
}

func (r *resp) Okf(ctx *gin.Context, s string, vs ...any) {
	ctx.JSON(http.StatusOK, fmt.Sprintf(s, vs...))
}

func (r *resp) OkI18n(ctx *gin.Context, key string, vs ...any) {
	ctx.JSON(http.StatusOK, fmt.Sprintf(i18n.MustGetMessage(key), vs...))
}

func (r *resp) ServerErr(ctx *gin.Context) {
	ctx.Status(http.StatusInternalServerError)
}

func (r *resp) ServerErrf(ctx *gin.Context, fmt string, vs ...any) {
	ctx.String(http.StatusInternalServerError, fmt, vs...)
}

func (r *resp) ServerI18n(ctx *gin.Context, key string, vs ...any) {
	ctx.String(http.StatusInternalServerError, i18n.MustGetMessage(key), vs...)
}

func (r *resp) NotFound(ctx *gin.Context) {
	ctx.Status(http.StatusNotFound)
}

func (r *resp) NotFoundf(ctx *gin.Context, fmt string, vs ...any) {
	ctx.String(http.StatusNotFound, fmt, vs...)
}

func (r *resp) NotFoundI18n(ctx *gin.Context, key string, vs ...any) {
	ctx.String(http.StatusNotFound, i18n.MustGetMessage(key), vs...)
}

func (r *resp) BadReq(ctx *gin.Context) {
	ctx.Status(http.StatusBadRequest)
}

func (r *resp) BadReqf(ctx *gin.Context, fmt string, vs ...any) {
	ctx.String(http.StatusBadRequest, fmt, vs...)
}

func (r *resp) BadReqI18n(ctx *gin.Context, key string, vs ...any) {
	ctx.String(http.StatusBadRequest, i18n.MustGetMessage(key), vs...)
}

func (r *resp) NoAuth(ctx *gin.Context) {
	ctx.Status(http.StatusUnauthorized)
}

func (r *resp) NoAuthf(ctx *gin.Context, fmt string, vs ...any) {
	ctx.String(http.StatusUnauthorized, fmt, vs...)
}

func (r *resp) NoAuthI18n(ctx *gin.Context, key string, vs ...any) {
	ctx.String(http.StatusUnauthorized, i18n.MustGetMessage(key), vs...)
}

func (r *resp) PreferError(ctx *gin.Context, err error) {
	if errorx.IsPreferred(err) {
		perr := err.(*errorx.PreferredError)
		ctx.String(perr.Code(), perr.Error())
	} else {
		ctx.String(http.StatusBadRequest, err.Error())
	}
}

func (r *resp) Error(ctx *gin.Context, code int, err error) {
	ctx.String(code, err.Error())
}
