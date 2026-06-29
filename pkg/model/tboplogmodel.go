package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TbOpLogModel = (*customTbOpLogModel)(nil)

type (
	// TbOpLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTbOpLogModel.
	TbOpLogModel interface {
		tbOpLogModel
		withSession(session sqlx.Session) TbOpLogModel
	}

	customTbOpLogModel struct {
		*defaultTbOpLogModel
	}
)

// NewTbOpLogModel returns a model for the database table.
func NewTbOpLogModel(conn sqlx.SqlConn) TbOpLogModel {
	return &customTbOpLogModel{
		defaultTbOpLogModel: newTbOpLogModel(conn),
	}
}

func (m *customTbOpLogModel) withSession(session sqlx.Session) TbOpLogModel {
	return NewTbOpLogModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customTbOpLogModel) ListAuditorId(ctx context.Context) ([]string, error) {
	auditorIds := make([]string, 0)
	listQuery := fmt.Sprintf("select user_id, count(user_id) from %s group by user_id", m.table)

	list := make([]struct {
		UserId string `json:"user_id"`
		Count  int    `json:"count"`
	}, 0)

	err2 := m.conn.QueryRowsPartialCtx(ctx, &list, listQuery)
	if err2 != nil {
		return auditorIds, err2
	}

	for _, v := range list {
		auditorIds = append(auditorIds, v.UserId)
	}
	logx.WithContext(ctx).Infof(">>>>>>>> auditorIds: %v", auditorIds)

	return auditorIds, nil
}

func (m *customTbOpLogModel) GetOpLogListByCondition(ctx context.Context, offset, limit, opTimeFrom, OpTimeTo int64, opUserId, auditType, sort, reserve string) (int64, []*TbOpLog, error) {
	list := make([]*TbOpLog, 0)
	listSelect := fmt.Sprintf("select * from %s ", m.table)
	listWhere := "where 1 = 1 "

	// 条件查询
	var args []any
	if auditType != "" {
		listWhere = listWhere + " and audit_type = ? "
		args = append(args, auditType)
	}
	if opUserId != "" {
		listWhere = listWhere + " and user_id = ? "
		args = append(args, opUserId)
	}
	if opTimeFrom != 0 {
		listWhere = listWhere + " and op_time >= ? "
		args = append(args, opTimeFrom)
	}
	if OpTimeTo != 0 {
		listWhere = listWhere + " and op_time <= ? "
		args = append(args, OpTimeTo)
	}

	// 排序
	listSort := ""
	if sort != "" {
		listSort = " order by op_time " + reserve + " "
	} else {
		listSort = " order by " + sort + " " + reserve + " "
	}

	// 分页
	listPage := fmt.Sprintf(" limit %d offset %d ", limit, offset)

	// 统计total
	countQuery := listSelect + listWhere + listSort
	err := m.conn.QueryRowsPartialCtx(ctx, &list, countQuery, args...)
	if err != nil {
		return 0, nil, err
	}
	total := int64(len(list))

	// 统计具体数据
	list = list[:0]
	listQuery := listSelect + listWhere + listSort + listPage
	err2 := m.conn.QueryRowsPartialCtx(ctx, &list, listQuery, args...)
	if err2 != nil {
		return 0, nil, err2
	}

	return total, list, nil
}

func (m *customTbOpLogModel) GetOpLogById(ctx context.Context, opId int64) (*TbOpLog, error) {

	list := make([]*TbOpLog, 0)

	querySelect := fmt.Sprintf("select * from %s ", m.table)

	// 条件查询
	var args []any
	queryWhere := " where 1 = 1 and id = ? "
	args = append(args, opId)

	listQuery := querySelect + queryWhere

	err := m.conn.QueryRowsPartialCtx(ctx, &list, listQuery, args...)
	if err != nil {
		return nil, err
	}

	return list[0], nil
}

func (m *customTbOpLogModel) Add(ctx context.Context, data *TbOpLog) (int64, error) {
	insertRes, err := m.Insert(ctx, data)
	if err != nil {
		return 0, err
	}
	return insertRes.LastInsertId()
}
