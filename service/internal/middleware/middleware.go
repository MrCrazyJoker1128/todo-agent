package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func MiddlewareResponse(r *ghttp.Request) {
	r.Middleware.Next()

	if r.Response.BufferLength() > 0 {
		return
	}

	err := r.GetError()
	if err != nil {
		r.Response.WriteJson(Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	data := r.GetHandlerResponse()
	if data == nil {
		r.Response.WriteJson(Response{
			Code:    200,
			Message: "success",
		})
	} else {
		r.Response.WriteJson(Response{
			Code:    200,
			Message: "success",
			Data:    data,
		})
	}
}
