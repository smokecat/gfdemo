package sorm

import (
  "fmt"
  "strings"

  "github.com/gogf/gf/frame/g"
  "github.com/gogf/gf/text/gstr"
)

const (
  OrderPosAsc  = "asc"
  OrderPosDesc = "desc"
)

type OrderBy struct {
  // OrderColumns should conform to `model,col1,col2` and case-sensitive.
  OrderColumns string `json:"orderColumns" dc:"逗号分割的model字段字符串,不能有空格"`

  // OrderPositions is a comma-separated string of asc or desc(d), which corresponds to the column name one by one.
  // If it is not equal to desc(d) or is an empty string, it will be set to asc by default.
  OrderPositions string `json:"orderPosition" dc:"逗号分割，对应column的排序"`
}

type order struct {
  Column   string
  Position string
}

// ToOrderByParam convert orderby to db params.
func (o OrderBy) ToOrderByParam() string {
  if o.OrderColumns == "" {
    return ""
  }

  columns := strings.Split(o.OrderColumns, ",")
  positions := strings.Split(string(o.OrderPositions), ",")

  var builder strings.Builder
  for i, v := range columns {
    var (
      pos = OrderPosAsc
    )
    if i < len(positions) && (positions[i] == OrderPosDesc || positions[i] == "d") {
      pos = OrderPosDesc
    }
    builder.WriteString(g.DB().GetCore().QuoteString(fmt.Sprintf("%s %s", gstr.CaseSnake(v), pos)))
    builder.WriteString(", ")
  }
  res := builder.String()[:builder.Len()-2]

  g.Log().Debugf("order by convert res: %v - %+v", res, o)
  return res
}

func (o OrderBy) ToOrders() []order {
  if o.OrderColumns == "" {
    return []order{}
  }

  columns := strings.Split(o.OrderColumns, ",")
  positions := strings.Split(string(o.OrderPositions), ",")
  orders := make([]order, len(columns))

  for i, v := range columns {
    var (
      pos = OrderPosAsc
    )
    if i < len(positions) && (positions[i] == OrderPosDesc || positions[i] == "d") {
      pos = OrderPosDesc
    }

    orders = append(orders, order{gstr.CaseSnake(v), pos})
  }
  return orders
}
