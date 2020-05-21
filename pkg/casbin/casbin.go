package casbin

import (
	"os"

	"github.com/qinyuanmao/go-utils/logutl"

	"github.com/casbin/casbin"
	. "github.com/casbin/gorm-adapter"
	_ "github.com/go-sql-driver/mysql"
)

//Casbin 权限控制
func Casbin() (*casbin.Enforcer, error) {
	CasbinAdapter := NewAdapter("mysql", os.Getenv("MYSQL_URL"))
	e := casbin.NewEnforcer("./conf/rbac_model.conf", CasbinAdapter)

	if err := e.LoadPolicy(); err == nil {
		return e, err
	} else {
		logutl.Debug("casbin rbac_model or policy init error, message: %v", err)
		return nil, err
	}
}
