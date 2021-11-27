package api

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"

	"github.com/smokecat/gfdemo/app/model"
	"github.com/smokecat/gfdemo/app/service"
	"github.com/smokecat/gfdemo/library/response"
)

var User = userApi{}

type userApi struct{}

func (u *userApi) Create(r *ghttp.Request) {
	var (
		apiReq    *model.UserApiCreateReq
		serviceIn *model.UserServiceCreateInput
	)

	if err := r.ParseForm(&apiReq); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	if err := gconv.Struct(apiReq, &serviceIn); err != nil {
		response.JsonExit(r, 1, err.Error())
	}

	if res, err := service.User.Create(r.Context(), serviceIn); err != nil {
		g.Log().Errorf("create user failed: %v", err)
		response.JsonExit(r, 1, err.Error())
	} else {
		response.JsonExit(r, 0, "success", res)
	}
}

func (u *userApi) List(r *ghttp.Request) {
	var (
		apiReq    *model.UserApiListReq
		serviceIn *model.UserServiceListInput
	)

	if err := r.ParseQuery(&apiReq); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	if err := gconv.Struct(apiReq, &serviceIn); err != nil {
		response.JsonExit(r, 1, err.Error())
	}

	if res, err := service.User.List(r.Context(), serviceIn); err != nil {
		g.Log().Errorf("create user failed: %v", err)
		response.JsonExit(r, 1, err.Error())
	} else {
		response.JsonExit(r, 0, "success", res)
	}
}
