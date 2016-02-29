package main

import (
	"os"
	"runtime"

	"github.com/codegangsta/cli"

	"github.com/kuuyee/gogs-learn/cmd"
	"github.com/kuuyee/gogs-learn/modules/setting"
)

const APP_VER = "KuuYee.0.8.25.0129"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU()) //设置使用的CPU核数为本机CPU核数
	setting.AppVer = APP_VER
}

func main() {
	app := cli.NewApp()
	app.Name = "Gogs-Learn"
	app.Usage = "Gogs项目源码学习"
	app.Version = APP_VER
	app.Commands = []cli.Command{
		cmd.CmdWeb,
	}

	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}
