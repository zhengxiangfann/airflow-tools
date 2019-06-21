package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.html"
}

func (c *MainController) Post() {
	DAG_PATH := "/Users/zhengxiangfan/code/go_src/1.12/src/github.com/upcode/static/upfile/"
	f, h, _ := c.GetFile("myfile")
	path := DAG_PATH + h.Filename
	f.Close()
	c.SaveToFile("myfile", path)
	c.TplName = "maincontroller/post.html"
}
