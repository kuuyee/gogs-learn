// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"io/ioutil"
	"path"

	"gopkg.in/ini.v1"
	"gopkg.in/macaron.v1"

	"github.com/codegangsta/cli"
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

	"github.com/kuuyee/gogs-learn/modules/log"
	"github.com/kuuyee/gogs-learn/modules/setting"
	"github.com/kuuyee/gogs-learn/modules/template"
	"github.com/kuuyee/gogs-learn/routers"
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
		{"github.com/go-xorm/xorm", func() string { return xorm.Version }, "0.5.2.0304"},
		{"github.com/go-macaron/binding", binding.Version, "0.2.1"},
		{"github.com/go-macaron/cache", cache.Version, "0.1.2"},
		{"github.com/go-macaron/csrf", csrf.Version, "0.1.0"},
		{"github.com/go-macaron/i18n", i18n.Version, "0.2.0"},
		{"github.com/go-macaron/session", session.Version, "0.1.6"},
		{"github.com/go-macaron/toolbox", toolbox.Version, "0.1.0"},
		{"gopkg.in/ini.v1", ini.Version, "1.8.4"},
		{"gopkg.in/macaron.v1", macaron.Version, "1.1.2"},
		{"github.com/gogits/git-module", git.Version, "0.2.9"},
		{"github.com/gogits/go-gogs-client", gogs.Version, "0.7.4"},
	}

	for _, c := range checkers {
		if !version.Compare(c.Version(), c.Expected, ">=") {
			log.Fatal(4, "Package '%s' version is too old (%s -> %s), did you forget to update?", c.ImportPath, c.Version(), c.Expected)
		}
	}
}

func newMacaron() *macaron.Macaron {
	m := macaron.New()

	// DISABLE_ROUTER_LOG: 激活该选项来禁止打印路由日志
	// 判断是否禁用，如果禁用则引入macaron日志
	if !setting.DisableRouterLog {
		m.Use(macaron.Logger())
	}
	// 引入macaron恢复机制
	m.Use(macaron.Recovery())

	if setting.Protocol == setting.FCGI {
		m.SetURLPrefix(setting.AppSubUrl)
	}

	// 设定静态资源路径
	m.Use(macaron.Static(
		path.Join(setting.StaticRootPath, "public"),
		macaron.StaticOptions{
			SkipLogging: setting.DisableRouterLog,
		},
	))
	m.Use(macaron.Static(
		setting.AvatarUploadPath,
		macaron.StaticOptions{
			Prefix:      "avatars",
			SkipLogging: setting.DisableRouterLog,
		},
	))

	// 设置渲染模板
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory:         path.Join(setting.StaticRootPath, "templates"),
		AppendDirectories: []string{path.Join(setting.CustomPath, "templates")},
		Funcs:             template.NewFuncMap(),
		IndentJSON:        macaron.Env != macaron.PROD,
	}))
	return m

}
func runWeb(ctx *cli.Context) {
	if ctx.IsSet("config") {
		setting.CustomConf = ctx.String("config")
	}
	routers.GlobalInit()
	checkVersion()
}
