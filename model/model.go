package po

import "time"

type Refund struct {
	Id           int64     `xorm:"bigint(20) unsigned pk  autoincr  notnull 'id'"` // 主键
	OrderId      int64     `xorm:"bigint(20) notnull 'order_id'"`                  // 冗余订单ID
	RefundAmount float64   `xorm:"decimal(20,4) notnull 'refund_amount'"`          // 申请退款金额
	FinalAmount  float64   `xorm:"decimal(20,4) notnull 'final_amount'"`           // 最终成功退款金额
	State        int64     `xorm:"int(11) notnull 'state'"`                        // 退款状态
	Description  string    `xorm:"varchar(128) notnull 'description'"`             // 退款描述
	RefundTime   time.Time `xorm:"datetime notnull 'refund_time'"`                 // 最终退款到账时间
	CreateTime   time.Time `xorm:"datetime notnull 'create_time'"`                 // 创建时间
	UpdateTime   time.Time `xorm:"datetime notnull 'update_time'"`                 // 更新时间
}

func (Refund) TableName() string {
	return "refund"
}
