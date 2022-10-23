package geek

import (
	"encoding/json"
	"io"
	"net/http"
)

type Context struct {
	W http.ResponseWriter
	R *http.Request
}

// ReadJson 读取body 反序列化成json
func (c *Context) ReadJson(req interface{}) error {
	r := c.R
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, req)
	if err != nil {
		return err
	}
	// golang 语言特点，一般直接改变数据，而不是返回数据
	return nil
}

// Write
func (c *Context) Write(code int, resp []byte) error {
	c.W.WriteHeader(code)
	_, err := c.W.Write(resp)
	return err
}

// WriteJson
func (c *Context) WriteJson(code int, resp interface{}) error {
	respJson, err := json.Marshal(resp)
	if err == nil {
		err = c.Write(code, respJson)
	}
	return err
}

func (c *Context) NotFound404() error {
	return c.Write(http.StatusNotFound, []byte("Not Found"))
}

func (c *Context) OkJson(resp interface{}) error {
	return c.WriteJson(http.StatusOK, resp)
}

func (c *Context) SystemErrJson(resp interface{}) error {
	return c.WriteJson(http.StatusInternalServerError, resp)
}

func (c *Context) BadRequestJson(resp interface{}) error {
	return c.WriteJson(http.StatusBadRequest, resp)
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W: w,
		R: r,
	}
}
