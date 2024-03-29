// Code generated by protoc-gen-go-hip. DO NOT EDIT.
// versions:v1.2.0

package gen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the hip package it is being compiled against.
import (
	"net/http"

	"github.com/GodWY/gutil"

	"github.com/gin-gonic/gin"
)

// generated http method
func registerLoginHttpHandler(srv *gin.Engine, h LoginHttpHandler) {
	group := srv.Group("api/xxx", gin.Logger(), gin.Recovery())
	group.GET("/getUserName", h.GetUserName)
	group.POST("/getUserID", gin.Logger(), gin.Recovery(), h.GetUserID, gin.Recovery())
}

var httpLogin Login

func RegisterLoginHttpHandler(srv *gin.Engine, h Login) {
	tmp := new(xLogin)
	registerLoginHttpHandler(srv, tmp)
	httpLogin = h
}

//  Login this is a test
type Login interface {
	// GetUserName
	GetUserName(ctx *gin.Context, in *Request) (out *Response, err error)
	// GetUserID
	GetUserID(ctx *gin.Context, in *Request) (out *Response, err error)
}

// generated http handle
type LoginHttpHandler interface {
	GetUserName(ctx *gin.Context)
	GetUserID(ctx *gin.Context)
}

type xLogin struct{}

func (x *xLogin) GetUserName(ctx *gin.Context) {
	req := &Request{}
	if err := ctx.ShouldBind(req); err != nil {
		detail := "bind request error: " + err.Error()
		rt := gutil.RetFail(10000, detail)
		ctx.JSON(http.StatusOK, rt)
		return
	}
	rsp, err := httpLogin.GetUserName(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, gutil.RetError(err))
		return
	}
	ctx.JSON(http.StatusOK, gutil.RetSuccess("success", rsp))
}

func (x *xLogin) GetUserID(ctx *gin.Context) {
	req := &Request{}
	if err := ctx.ShouldBind(req); err != nil {
		detail := "bind request error: " + err.Error()
		rt := gutil.RetFail(10000, detail)
		ctx.JSON(http.StatusOK, rt)
		return
	}
	rsp, err := httpLogin.GetUserID(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, gutil.RetError(err))
		return
	}
	ctx.JSON(http.StatusOK, gutil.RetSuccess("success", rsp))
}
