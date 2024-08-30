package user

// User is a struct for table user
//
//db:table user
type User struct {
	UserID   int    `db:"id"  pk:"true"`
	Username string `db:"name"`
	GroupID  int    `db:"group_id"`
}

type UserSlice []*User

// AdminUser is a struct for table admin_user
//
//db:table admin_user
type AdminUser struct {
	AdminUserID int    `db:"id"  pk:"true"`
	Username    string `db:"name"`
	GroupID     int    `db:"group_id"`
}
