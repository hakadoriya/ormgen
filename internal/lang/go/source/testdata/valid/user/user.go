package user

// User is a struct for table user
//
//db:table user
type User struct {
	UserID   int    `db:"user_id"  pk:"true"`
	Username string `db:"username"           hasOne:"Username"  hasMany:"UsernameAndAddress"`
	Address  string `db:"address"                               hasMany:"UsernameAndAddress"`
	GroupID  int    `db:"group_id"`
}

type UserSlice []*User

// AdminUser is a struct for table admin_user
//
//db:table admin_user
type AdminUser struct {
	AdminUserID int    `db:"admin_user_id"  pk:"true"`
	Username    string `db:"username"`
	GroupID     int    `db:"group_id"`
}

type AdminUserSlice []*AdminUser

// InvalidStruct is a struct but has no db tag for field
//
//db:table noop
type InvalidStruct struct {
	InvalidID int // has no db tag
}
