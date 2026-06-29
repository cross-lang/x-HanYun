package model

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var XXXDBConn sqlx.SqlConn
var TbOpLogIns TbOpLogModel

func Init(XXXMySQLDsn string) {
	if XXXDBConn != nil {
		return
	}

	XXXDBConn = sqlx.NewMysql(XXXMySQLDsn)
	if err := checkConn(XXXDBConn); err != nil {
		panic(err)
	}

	TbOpLogIns = NewTbOpLogModel(XXXDBConn)
}

// 检查连接是否有效
func checkConn(conn sqlx.SqlConn) error {
	var tables []string
	err := conn.QueryRows(&tables, "SHOW TABLES")
	if err != nil {
		return err
	}
	return nil
}
