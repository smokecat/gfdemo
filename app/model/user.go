package model

import (
	"github.com/smokecat/gfdemo/library/sorm"
)

type UserApiCreateReq struct {
	Name     string `json:"name" v:"required#缺少name"`
	Nick     string `json:"nick" v:"required#缺少nick"`
	Password string `json:"password" v:"required#缺少password"`
	Email    string `json:"email"`
	Age      uint   `json:"age"`
	HeadImg  string `json:"head_img"`
}

type UserApiListReq struct {
	sorm.Paging
	sorm.OrderBy `v:"order-by:user,id,age,headImg#"`
}

type UserServiceCreateInput struct {
	Name     string
	Nick     string
	Password string
	Email    string
	Age      uint
	HeadImg  string
}

type UserServiceCreateOutput struct {
	Id uint64
}

type UserServiceListInput struct {
	sorm.Paging
	sorm.OrderBy
}

type UserServiceListOutput struct {
	Users []*User
}
