package main

import (
	"arxiv/database"
	_ "arxiv/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	database.InitDB()
	beego.Run("0.0.0.0:8083")

}
