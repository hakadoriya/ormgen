package user

// User is a struct for table user
//
//db:table user
type User struct {
	ID   int    `db:"id"  pk:"true"`
	Name string `db:"name"`
}

type UserSlice []*User
