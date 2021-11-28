package boot

import (
  "github.com/gogf/gf/frame/g"
  "github.com/gogf/gf/util/gvalid"

  "github.com/smokecat/gfdemo/app/model"
  "github.com/smokecat/gfdemo/library/sorm"
  _ "github.com/smokecat/gfdemo/packed"
)

func init() {
  // Register sorm model columns map
  sorm.RegisterModelColumnsMap(map[string]interface{}{"user": model.User{}})

  // Register order-by rule
  if err := gvalid.RegisterRule(sorm.OrderByRuleName, sorm.RuleOrderBy); err != nil {
    g.Log().Panicf("register rule `%s` failed", sorm.OrderByRuleName)
  }
}
