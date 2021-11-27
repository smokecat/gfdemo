package boot

import (
  "github.com/gogf/gf/frame/g"
  "github.com/gogf/gf/util/gvalid"

  "github.com/smokecat/gfdemo/app/model"
  "github.com/smokecat/gfdemo/library/sorm"
  _ "github.com/smokecat/gfdemo/packed"
)

func init() {
  // Register sorm model fields map
  sorm.RegisterModelFieldsMap(map[string]interface{}{"user": model.User{}})

  // Register custom rule
  if err := gvalid.RegisterRule("contain-model-fields", sorm.RuleContainModelFields); err != nil {
    g.Log().Panic("register rule contain-model-fields failed")
  }
}
