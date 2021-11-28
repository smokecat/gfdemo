# A GoFrame Project Demo

> [GoFrame](https://goframe.org)

## package sorm

数据库相关功能实现

### paging 分页

自定义了一个包含分页参数的结构体，并添加了gvalid校验规则

```go
package sorm

type Paging struct {
  Page int `json:"page" v:"required|integer|min:1#||"`
  Size int `json:"size" v:"required|integer|min:1#||"`
}
```

在定义请求接口时，可添加Paging到结构体，作为匿名属性，Page和Size字段可提升为接口结构体的字段。

```go
package model

import (
  "github.com/smokecat/gfdemo/library/sorm"
)

type UserApiListReq struct {
  sorm.Paging
}
```

接口请求须携带page和size参数。

```http request
GET /api/user/list?page=1&size=10
```

可正常对参数进行校验

```go
package api

import (
  "github.com/gogf/gf/net/ghttp"
)

func List(r *ghttp.Request) {
  var (
    apiReq UserApiListReq
  )

  // apiReq = UserAoiListReq{Page: 1, Size: 10}
  if err := r.ParseQuery(&apiReq); err != nil {
    response.JsonExit(r, 1, err.Error())
  }
}
```

提供转换为Limit方法参数的方法: `func (p Paging) ToLimitParam() (offset int, limit int)`

```go
package service

import (
  "context"

  "github.com/smokecat/gfdemo/app/dao"
  "github.com/smokecat/gfdemo/app/model"
)

func List(ctx context.Context, in *model.UserServiceListInput) (*model.UserServiceListOutput, error) {
  var users []*model.User

  if err := dao.User.Ctx(ctx).Limit(in.ToLimitParam()).Scan(&users); err != nil {
    return nil, err
  }
  return &model.UserServiceListOutput{Users: users}, nil
}
```

### order by 排序

自定义了一个包含排序参数的结构体，并自定义了一个用户gvalid校验的结构体规则。

```go
package sorm

type OrderBy struct {
  // OrderColumns should conform to `model,col1,col2` and case-sensitive.
  OrderColumns string `json:"orderColumns" dc:"逗号分割的model字段字符串,不能有空格"`

  // OrderPositions is a comma-separated string of asc or desc(d), which corresponds to the column name one by one.
  // If it is not equal to desc(d) or is an empty string, it will be set to asc by default.
  OrderPositions string `json:"orderPosition" dc:"逗号分割，对应column的排序"`
}
```

在定义请求接口时，可添加OrderBy到结构体，作为匿名属性，OrderColumns和OrderPositions可提升为接口结构体的字段。

```go
package model

import (
  "github.com/smokecat/gfdemo/library/sorm"
)

type UserApiListReq struct {
  sorm.OrderBy `v:"order-by:user,id,age,headImg#"`
}
```

接口请求可携带orderColumns和orderPositions参数（可选）。

```http request
GET /api/user/list?orderColumns=id,age&orderPositions=asc,desc
```

校验OrderBy参数: <a href="#order-by">详见order-by规则</a>

提供转换为OrderBy方法参数的方法: `func (o OrderBy) ToOrderByParam() string`

```go
package service

func (u *userService) List(ctx context.Context, in *model.UserServiceListInput) (*model.UserServiceListOutput, error) {
  var users []*model.User

  // Order("id asc, age desc")
  if err := dao.User.Ctx(ctx).Order(in.ToOrderByParam()).Scan(&users); err != nil {
    return nil, err
  }
  return &model.UserServiceListOutput{Users: users}, nil
}
```

ToOrderByParam方法转换规则
- OrderPostions逗号分割的排序规则与OrderColumns一一对应。
- 如果对应的排序方向不等于"desc"、"d"或者长度不足，那么将转换为"asc"。

#### <a id="order-by"></a> OrderBy rule

为了更加通用，OrderBy字段参数的校验较为复杂。推荐按如下步骤进行。

1. 在`boot`包的init方法中注册数据库model。并且全局注册order-by校验规则。

```go
package boot

func init() {
  // Register sorm model columns map
  sorm.RegisterModelColumnsMap(map[string]interface{}{"user": model.User{}})

  // Register order-by rule
  if err := gvalid.RegisterRule(sorm.OrderByRuleName, sorm.RuleOrderBy); err != nil {
    g.Log().Panicf("register rule `%s` failed", sorm.OrderByRuleName)
  }

}
```

2. 在OrderBy匿名字段后添加gvalid校验标签`order-by`。
```go
package model

type UserApiListReq struct {
  sorm.OrderBy `v:"order-by:user,id,age,headImg#"`
}
```

标签校验规则如下：
  - 形如`order-by:<model>,<col1>,<col2>...`。
  - 规则将逗号分割的第一个元素作为model，后续元素作为列名，并且都是可选的。
  - 如果未指定model，即只有规则名`order-by`，那么默认允许所有参数。
  - 如果仅包含model，即`order-by:model`，那么有以下两种情况。
    - 该model在步骤1中已注册，那么值允许参数包含注册结构体(model.User)的所有列(包含列名的小驼峰、大驼峰、下划线格式)。
    - 若为注册，那么默认允许所有参数。
  - 如果包含cols，那么无论指定的model是否注册，只允许参数包含cols中的列名(大小写敏感)。

3. 使用gvalid进行参数校验
```go
package api

func List(r *ghttp.Request) {
  var (
    apiReq UserApiListReq
  )

  // apiReq = UserAoiListReq{OrderColumns: "id,age", OrderPositions: "asc,desc"}
  if err := r.ParseQuery(&apiReq); err != nil {
    response.JsonExit(r, 1, err.Error())
  }
}
```