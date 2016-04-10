package models

import "github.com/go-xorm/xorm"

var (
	x *xorm.Engine
)

func init() {

}

func Ping() error {
	return x.Ping()
}
