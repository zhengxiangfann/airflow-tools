package routers

import (
	"github.com/upcode/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("/dagdraw", &controllers.DagController{})
	beego.Router("/dagdata", &controllers.DagController{}, "*:Dag")
}
