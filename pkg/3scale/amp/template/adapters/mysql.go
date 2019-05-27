package adapters

import (
	"github.com/3scale/3scale-operator/pkg/3scale/amp/component"
	templatev1 "github.com/openshift/api/template/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Mysql struct {
}

func NewMysqlAdapter(options []string) Adapter {
	return NewAppenderAdapter(&Mysql{})
}

func (m *Mysql) Parameters() []templatev1.Parameter {
	return []templatev1.Parameter{
		templatev1.Parameter{
			Name:        "MYSQL_USER",
			DisplayName: "MySQL User",
			Description: "Username for MySQL user that will be used for accessing the database.",
			Value:       "mysql",
			Required:    true,
		},
		templatev1.Parameter{
			Name:        "MYSQL_PASSWORD",
			DisplayName: "MySQL Password",
			Description: "Password for the MySQL user.",
			Generate:    "expression",
			From:        "[a-z0-9]{8}",
			Required:    true,
		},
		templatev1.Parameter{
			Name:        "MYSQL_DATABASE",
			DisplayName: "MySQL Database Name",
			Description: "Name of the MySQL database accessed.",
			Value:       "system",
			Required:    true,
		},
		templatev1.Parameter{
			Name:        "MYSQL_ROOT_PASSWORD",
			DisplayName: "MySQL Root password.",
			Description: "Password for Root user.",
			Generate:    "expression",
			From:        "[a-z0-9]{8}",
			Required:    true,
		},
	}
}

func (m *Mysql) Objects() ([]runtime.RawExtension, error) {
	mysqlOptions, err := m.options()
	if err != nil {
		return nil, err
	}
	mysqlComponent := component.NewMysql(mysqlOptions)
	return mysqlComponent.Objects(), nil
}

func (a *Mysql) options() (*component.MysqlOptions, error) {
	mob := component.MysqlOptionsBuilder{}
	mob.AppLabel("${APP_LABEL}")
	mob.DatabaseName("${MYSQL_DATABASE}")
	mob.User("${MYSQL_USER}")
	mob.Password("${MYSQL_PASSWORD}")
	mob.RootPassword("${MYSQL_ROOT_PASSWORD}")
	mob.DatabaseURL("mysql2://root:" + "${MYSQL_ROOT_PASSWORD}" + "@system-mysql/" + "${MYSQL_DATABASE}")
	return mob.Build()
}
