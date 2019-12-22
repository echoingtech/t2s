package po

import "time"

type UserSyncs struct {
	UserId    int64     `xorm:"bigint(20) pk  notnull 'user_id'"` // 用户 ID
	Synced    int       `xorm:"tinyint(1) null 'synced'"`         // 是否已同步
	CreatedAt time.Time `xorm:"datetime notnull 'created_at'"`
	UpdatedAt time.Time `xorm:"datetime notnull 'updated_at'"`
}

func (UserSyncs) TableName() string {
	return "user_syncs"
}
