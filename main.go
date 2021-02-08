package main

import (
	"Blog/route"

	"Blog/model"
	"Blog/util"
	_ "Blog/util/validator"
)

func main() {
	db := model.GetDb()
	defer db.Close()
	route.InitRoute().Run(util.HttpPort)
}
