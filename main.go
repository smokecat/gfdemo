package main

import (
	_ "github.com/smokecat/gfdemo/boot"
	_ "github.com/smokecat/gfdemo/router"

	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Server().Run()
}
