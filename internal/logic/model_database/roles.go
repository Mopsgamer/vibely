package model_database

type Role struct {
	GroupId uint64 `db:"group_id"`
	UserId  uint64 `db:"user_id"`
	RightId uint32 `db:"right_id"`
}
