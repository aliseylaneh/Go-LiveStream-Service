package models

type Pagination struct {
	Offset   int32
	Limit    int32
	GetTotal bool
}
