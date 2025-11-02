package routers

import (
	"arxiv/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.AuthController{}, "get:GetLogin")

	beego.Router("/login", &controllers.AuthController{}, "get:GetLogin;post:PostLogin")

	beego.Router("/register", &controllers.AuthController{}, "get:GetRegister;post:PostRegister")

	beego.Router("/dashboard/:id", &controllers.DashboardController{}, "get:Get;post:Post")
	// routers/router.go

	beego.Router("/note/toggle/:id", &controllers.NoteController{}, "post:Toggle")

	beego.Router("/note/delete/:id", &controllers.NoteController{}, "post:Delete")

	beego.Router("/dashboard/:id", &controllers.DashboardController{}, "get:Get")

	beego.Router("/logout", &controllers.AuthController{}, "get:Logout")

	beego.Router("/help", &controllers.RegisterController{}, "get:Help;post:Help")

	beego.Router("/Buyurtma", &controllers.RegisterController{}, "get:Buyurtma;post:Buyurtma")

}
