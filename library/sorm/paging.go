package sorm

type Paging struct {
	Page int `json:"page" v:"required|integer|min:1#||"`
	Size int `json:"size" v:"required|integer|min:1#||"`
}

// ToLimitParam convert page and size to offeset and take
func (p Paging) ToLimitParam() (int, int) {
	return (p.Page - 1) * p.Size, p.Size
}

func NewPaging() Paging {
	return Paging{}
}