// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/kuuyee/gogs-learn/modules/log"
	"github.com/kuuyee/gogs-learn/modules/setting"
	"github.com/kuuyee/gogs-learn/routers"

	_ "github.com/Unknwon/log"
	"github.com/Unknwon/macaron"
	"github.com/codegangsta/cli"
	"github.com/go-ini/ini"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/i18n"
	"github.com/go-macaron/session"
	"github.com/go-macaron/toolbox"
	"github.com/go-xorm/xorm"
	"github.com/gogits/git-module"
	"github.com/gogits/go-gogs-client"
	"github.com/mcuadros/go-version"
	"io/ioutil"
)

var CmdWeb = cli.Command{
	Name:  "web",
	Usage: "启动Gogs web服务器",
	Description: `Gogs web server is the only thing you need to run,
and it takes care of all the other things for you`,
	Action: runWeb,
	Flags: []cli.Flag{
		stringFlag("port, p", "3000", "Temporary port number to prevent conflict"),
		stringFlag("config, c", "custom/conf/app.ini", "Custom configuration file path"),
	},
}

type VerChecker struct {
	ImportPath string
	Version    func() string
	Expected   string
}

func runWeb(ctx *cli.Context) {
	if ctx.IsSet("config") {
		setting.CustomConf = ctx.String("config")
	}
	routers.GlobalInit()
}

// checkVersion checks if binary matches the version of templates files.
func checkVersion() {
	// Templates.
	data, err := ioutil.ReadFile(setting.StaticRootPath + "/templates/.VERSION")
	if err != nil {
		log.Fatal(4, "Fail to read 'templates/.VERSION': %v", err)
	}
	if string(data) != setting.AppVer {
		//log.Fatal(4, "Binary和templates文件版本不匹配，你是不是忘了重新编译?")
	}

	// 检查依赖版本
	checkers := []VerChecker{
		{"github.com/go-xorm/xorm", func() string { return xorm.Version }, "0.4.4.1029"},
		{"github.com/go-macaron/binding", binding.Version, "0.1.0"},
		{"github.com/go-macaron/cache", cache.Version, "0.1.2"},
		{"github.com/go-macaron/csrf", csrf.Version, "0.0.3"},
		{"github.com/go-macaron/i18n", i18n.Version, "0.2.0"},
		{"github.com/go-macaron/session", session.Version, "0.1.6"},
		{"github.com/go-macaron/toolbox", toolbox.Version, "0.1.0"},
		{"gopkg.in/ini.v1", ini.Version, "1.8.4"},
		{"gopkg.in/macaron.v1", macaron.Version, "0.8.0"},
		{"github.com/gogits/git-module", git.Version, "0.2.4"},
		{"github.com/gogits/go-gogs-client", gogs.Version, "0.7.2"},
	}

	for _, c := range checkers {
		if !version.Compare(c.Version(), c.Expected, ">=") {
			//log.Fatal(4, "Package '%s' version is too old (%s -> %s), did you forget to update?", c.ImportPath, c.Version(), c.Expected)
		}
	}
}
