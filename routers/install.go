// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package routers

import (
	"github.com/kuuyee/gogs-learn/modules/base"
	"github.com/kuuyee/gogs-learn/modules/setting"
)

const (
	INSTALL base.TplName = "install"
)

// GlobalInit is for global configuration reload-able.
// GlobalInit 用来全局配置,可重载
func GlobalInit() {
	//初始化配置
	setting.NewContext()
}
