package routers

import (
	"arxiv/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.AuthController{}, "get:GetLogin")

	beego.Router("/login", &controllers.AuthController{}, "get:GetLogin;post:PostLogin")

	beego.Router("/register", &controllers.RegisterController{}, "get:Get;post:Post")

	beego.Router("/dashboard/:id", &controllers.DashboardController{}, "get:Get;post:Post")
	// routers/router.go
	beego.Router("/admin", &controllers.AdminController{}, "get:Get;post:Post")
	beego.Router("/admin/login", &controllers.AdminController{}, "get:Login;post:LoginPost")
	beego.Router("/admin/logout", &controllers.AdminController{}, "get:Logout")
	beego.Router("/admin/add", &controllers.AdminController{}, "post:Add")

	beego.Router("/admin/create", &controllers.AdminController{}, "post:Post")
	beego.Router("/admin/delete", &controllers.AdminController{}, "get:Delete")

	// login sahifalari:
	beego.Router("/admin/login", &controllers.AdminController{}, "get:Login;post:LoginPost")

	beego.Router("/note/toggle/:id", &controllers.NoteController{}, "post:Toggle")

	beego.Router("/note/delete/:id", &controllers.NoteController{}, "post:Delete")

	beego.Router("/dashboard/:id", &controllers.DashboardController{}, "get:Get")

	beego.Router("/logout", &controllers.AuthController{}, "get:Logout")

	beego.Router("/help", &controllers.RegisterController{}, "get:Help;post:HelpPost")

	beego.Router("/Buyurtma", &controllers.RegisterController{}, "get:Buyurtma;post:BuyurtmaPost")

}
