package main

import (
	"os"
	"study-spider/boss"
	"study-spider/cnblogs"
	"study-spider/config"
	_ "study-spider/mail"
)

func main() {
	runmode := os.Getenv("runmode")
	if runmode == "" {
		runmode = "dev"
	}
	config.Init(runmode)
	// 运行博客园的爬虫
	blog := cnblogs.NewBlog()
	go blog.Run()
	// 运行boss爬虫
	bs := boss.NewBoss()
	bs.Run()
	select {}
}
