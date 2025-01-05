package group

// Group is a struct for table group
//
//db:table group
type Group struct {
	ID   int    `db:"id"  pk:"true"`
	Name string `db:"name"`
}

type GroupSlice []*Group
