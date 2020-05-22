package casbin

import (
	"os"

	"github.com/casbin/casbin"
	. "github.com/casbin/gorm-adapter"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

//Casbin 权限控制
func Casbin() (*casbin.Enforcer, error) {
	CasbinAdapter := NewAdapter("mysql", os.Getenv("MYSQL_URL"))
	e := casbin.NewEnforcer("./conf/rbac_model.conf", CasbinAdapter)

	if err := e.LoadPolicy(); err == nil {
		return e, err
	} else {
		logrus.Debug("casbin rbac_model or policy init error, message: %v", err)
		return nil, err
	}
}
