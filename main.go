package main

import (
	_ "github.com/upcode/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.SetStaticPath("/", "static")
	beego.Run()
}

