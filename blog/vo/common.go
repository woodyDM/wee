package vo

import (
	"time"
	"wee-server/support/database"
	"wee-server/wee"
)

var Zone08 = time.FixedZone("CST", 8*3600)

type Response struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Detail  string      `json:"detail"`
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}

func NewResponse(data interface{}) *Response {
	return &Response{
		Data:    data,
		Success: true,
	}
}

func NewFailResponse(msg string) *Response {
	return &Response{
		Code: 1,
		Msg:  msg,
	}
}

func WriteToResponse(ctx *wee.Context, msg string) {
	resp := NewFailResponse(msg)
	ctx.Json(resp)
}

type PageResp struct {
	Content []interface{} `json:"content"`
	Page    int           `json:"page"`
	Size    int           `json:"size"`
	Total   int64         `json:"total"`
}

func ToPageResp(p *database.Page) PageResp {
	return PageResp{
		Content: p.Data,
		Page:    p.Page,
		Size:    p.PageSize,
		Total:   p.TotalElement,
	}
}

func GetZone8UnixTime() int {
	return int(time.Now().In(Zone08).Unix())
}
