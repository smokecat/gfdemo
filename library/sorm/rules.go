package sorm

import (
  "context"
  "errors"
  "reflect"
  "strings"

  "github.com/gogf/gf/container/garray"
  "github.com/gogf/gf/frame/g"
  "github.com/gogf/gf/text/gstr"
  "github.com/gogf/gf/util/gconv"
)

const (
  OrderByRuleName = "order-by"
)

var modelColumnsMap = make(map[string][]string)

// RegisterModelColumnsMap register model name map model columns.
// Be careful that modelColumnsMap will be overwritten while execute.
// This function should not be executed more than once.
// It is recommended to register in the init function of the boot package.
// This function will register the model column's CamelCase, camelCase, and snake-case forms.
func RegisterModelColumnsMap(m map[string]interface{}) {
  g.Log().Debugf("Register model columns map: %+v", m)
  modelColumnsMap = make(map[string][]string, len(m))
  for k, v := range m {
    model := reflect.TypeOf(v)
    fieldStruct := reflect.VisibleFields(model)
    columns := garray.NewSortedStrArray()
    columns.SetUnique(true)

    for _, v := range fieldStruct {
      // columns[i] = v.Name
      columns.Add(v.Name, gstr.CaseCamel(v.Name), gstr.CaseCamelLower(v.Name), gstr.CaseSnake(v.Name))
    }

    modelColumnsMap[k] = columns.Slice()
  }

  g.Log().Debugf("modelColumnsMap: %v", modelColumnsMap)
}

// RuleOrderBy check the OrderColumns and OrderPositions fields for a given data.
// The format of `rule` should conform to `order-by:model,col1,col2`. Note that cols are case-sensitive.
// If no cols or given model have not registered by RegisterModelColumnsMap function, default allow all columns.
//
// Detail rule to see checkOrderColumns and checkOrderPositions
func RuleOrderBy(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
  const (
    defaultMessage = "Contain invalid column: `:col`"
  )

  g.Log().Debugf("Check rule `%s`: rule = `%v` value = `%v` message = `%v` data = `%+v`", OrderByRuleName, rule, value, message,
    data)
  // Set default message
  if message == "The :attribute value is invalid" {
    message = defaultMessage
  }

  var (
    ruleSplits  = strings.Split(rule, ":")
    model       = ""
    ruleColumns []string
  )

  // Append ruleColumns
  if len(ruleSplits) == 2 {
    colSplits := strings.Split(ruleSplits[1], ",")
    model = colSplits[0]
    if len(colSplits) > 1 {
      ruleColumns = colSplits[1:]
    }
  }

  var orderBy = OrderBy{}
  gconv.Struct(data, &orderBy)
  g.Log()
  if err := gconv.Struct(data, orderBy); err != nil {
    // Return nil if empty
    if orderBy.OrderColumns == "" {
      return nil
    }

    // Check OrderColumns. Get OrderColumns field from data
    if err := checkOrderColumns(model, ruleColumns, orderBy.OrderColumns, message); err != nil {
      return err
    }

    // Check OrderPositions
    if err := checkOrderPositions(); err != nil {
      return err
    }

    return nil
  }

  g.Log().Debugf("result of `%s` resolution: model = `%v` ruleColumns = `%v`", OrderByRuleName, model, ruleColumns)

  return errors.New("given data can not convert to OrderBy struct")
}

// checkOrderColumns check OrderColumns field.
// check `value` is a comma-separated string that contains model columns.
//
// The message support `:col` variable and will be replaced by invalid column name.
//
// The format of `value` should conform to `col1,col2`.
func checkOrderColumns(model string, columns []string, value string, message string) error {
  var (
    modelCols = garray.NewStrArrayFrom(columns)
  )

  // Check model columns by model name
  if modelCols.Len() == 0 {
    // Default set all model columns by model name if no specify columns.
    // Golang do not support dynamic load class so if model name do not registered, allow all columns
    if v, ok := modelColumnsMap[model]; ok {
      modelCols.SetArray(v)
    }
  }

  orderCols := strings.Split(value, ",")
  for _, orderCol := range orderCols {
    if modelCols.Len() > 0 && !modelCols.Contains(orderCol) {
      g.Log().Debugf("Rule `%s` expect model columns: `%s` but get `%s`", OrderByRuleName, modelCols.Join(","),
        orderCol)
      return errors.New(strings.Replace(message, ":col", orderCol, -1))
    }
  }
  return nil
}

// checkOrderPositions check OrderPositions field.
func checkOrderPositions() error {
  return nil
}
