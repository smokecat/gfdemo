package sorm

type Paging struct {
	Page int `json:"page" v:"required|integer|min:1#||"`
	Size int `json:"size" v:"required|integer|min:1#||"`
}

// ToLimitParam convert page and size to offset and limit
func (p Paging) ToLimitParam() (offset int, limit int) {
	return (p.Page - 1) * p.Size, p.Size
}

func NewPaging() Paging {
	return Paging{}
}