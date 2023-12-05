package gzy

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]any

type Context struct {
	// http 对象
	Writer http.ResponseWriter
	Req    *http.Request
	// 请求信息
	Path   string
	Method string
	Params map[string]string
	// 响应码
	StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// PostForm 对 c.Req.FormValue 的一层封装,获取表单的参数
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 对 c.Req.URL.Query().Get 的一层封装,获取uri的参数
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 对 http.ResponseWriter.Header().Set 的一层封装,设置响应头信息
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// String 响应返回格式化字符串
func (c *Context) String(code int, format string, values ...any) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	_, _ = c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// Json 响应返回Json
func (c *Context) Json(code int, obj any) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	_, _ = c.Writer.Write(data)
}

// HTML 响应返回HTML网页
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	_, _ = c.Writer.Write([]byte(html))
}
