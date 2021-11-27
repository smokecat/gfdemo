package service

import (
  "context"

  "github.com/gogf/gf/util/gconv"

  "github.com/smokecat/gfdemo/app/dao"
  "github.com/smokecat/gfdemo/app/model"
)

var User = userService{}

type userService struct{}

func (u *userService) Create(ctx context.Context, in *model.UserServiceCreateInput) (*model.UserServiceCreateOutput,
    error) {
  var user *model.User
  if err := gconv.Struct(in, &user); err != nil {
    return nil, err
  }

  id, err := dao.User.Ctx(ctx).OmitEmpty().InsertAndGetId(user)
  return &model.UserServiceCreateOutput{Id: uint64(id)}, err
}

func (u *userService) List(ctx context.Context, in *model.UserServiceListInput) (*model.UserServiceListOutput, error) {
  var users []*model.User

  if err := dao.User.Ctx(ctx).Limit(in.ToLimitParam()).Order(in.
    ToOrderByParam()).Scan(&users); err != nil {
    return nil, err
  } else {
    return &model.UserServiceListOutput{Users: users}, nil
  }
}
