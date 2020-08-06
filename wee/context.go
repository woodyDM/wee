package wee

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Context struct {
	Response     http.ResponseWriter
	Request      *http.Request
	Session      Session
	PathVariable map[string]string
	Server       *Server
	parsedForm   bool
}

func (ctx *Context) GetPathVariable(key string) string {
	v, ok := ctx.PathVariable[key]
	if !ok {
		panic(fmt.Errorf("PathVariable %s not found. ", key))
	}
	return v
}
func (ctx *Context) GetRequestForm(key string, defaultValue string) string {
	v, ok := ctx.GetRequestForm0(key)
	if !ok {
		return defaultValue
	} else {
		return v
	}
}

func (ctx *Context) GetRequestFormInt(key string, defaultValue int) int {
	v, ok := ctx.GetRequestForm0(key)
	if !ok {
		return defaultValue
	} else {
		return toInt(v)
	}
}

func (ctx *Context) GetRequestForm0(key string) (string, bool) {
	if !ctx.parsedForm {
		err := ctx.Request.ParseForm()
		if err != nil {
			panic(err)
		}
		ctx.parsedForm = true
	}
	values := ctx.Request.Form[key]
	if len(values) == 0 {
		return "", false
	} else {
		return values[0], true
	}
}
func (ctx *Context) GetPathVariableInt(key string) int {
	return toInt(ctx.GetPathVariable(key))
}

func toInt(s string) int {
	r, e := strconv.Atoi(s)
	if e != nil {
		panic(e)
	}
	return r
}

func (ctx *Context) String(v string) {
	ctx.Response.Write([]byte(v))
}

func (ctx *Context) Sprintf(v string, p ...interface{}) {
	result := fmt.Sprintf(v, p...)
	ctx.Response.Write([]byte(result))
}

func (ctx *Context) Json(i interface{}) {
	bytes, e := json.Marshal(i)
	if e != nil {
		panic("Unable to json")
	}
	ctx.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.Response.Write(bytes)
}

func (ctx *Context) JsonRequest(r interface{}) {
	body := ctx.Request.Body
	defer body.Close()
	bytes, e := ioutil.ReadAll(body)
	if e != nil {
		panic(e)
	}
	e = json.Unmarshal(bytes, r)
	if e != nil {
		panic(e)
	}

}
