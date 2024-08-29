package user

// Other is a struct for table other
//
//db:table other
type Other struct {
	OtherID int `db:"other_id" pk:"true"`
}

type OtherSlice []*Other
