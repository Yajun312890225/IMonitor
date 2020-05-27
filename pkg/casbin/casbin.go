package casbin

import (
	"os"

	"github.com/casbin/casbin"
	. "github.com/casbin/gorm-adapter"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type Casbin struct {
	CasbinAdapter *Adapter
	Enforce       *casbin.Enforcer
}

var mycasbin *Casbin

// GetCasbin 获取Casbin对象
func GetCasbin() *Casbin {
	if mycasbin == nil {
		mycasbin = new(Casbin)
		mycasbin.CasbinAdapter = NewAdapter("mysql", os.Getenv("MYSQL_URL"))
		mycasbin.Enforce = casbin.NewEnforcer("./conf/rbac_model.conf", mycasbin.CasbinAdapter)
		if err := mycasbin.Enforce.LoadPolicy(); err != nil {
			logrus.Debug("casbin rbac_model or policy init error, message: ", err)
		}
	}
	return mycasbin
}
