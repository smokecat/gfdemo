package sorm

import (
  "context"
  "errors"
  "reflect"
  "strings"

  "github.com/gogf/gf/container/garray"
  "github.com/gogf/gf/frame/g"
)

var modelFieldsMap = make(map[string][]string)

// RegisterModelFieldsMap register model name map model fields.
// Be careful that modelFieldsMap will be overwritten while execute.
// This function should not be executed more than once.
// It is recommended to register in the init function of the boot package.
func RegisterModelFieldsMap(m map[string]interface{}) {
  g.Log().Debugf("Register model fields map: %+v", m)
  modelFieldsMap = make(map[string][]string, len(m))
  for k, v := range m {
    model := reflect.TypeOf(v)
    fieldStruct := reflect.VisibleFields(model)
    fields := make([]string, len(fieldStruct))

    for i, v := range fieldStruct {
      fields[i] = v.Name
    }

    modelFieldsMap[k] = fields
  }

  g.Log().Debugf("modelFieldsMap: %v", modelFieldsMap)
}

// RuleContainModelFields check `value` is a comma-separated string that contains model fields.
//
// The format of `rule` should conform to `contain-model-fields:model,field1,field2`.
// Note that fields are case-sensitive.
// If no fields or given model have not registered by RuleContainModelFieldsModelRegister function,
// default allow all fields.
//
// The message support `:field` variable and will be replaced by invalid field name.
//
// The format of `value` should conform to `field1,field2`.
func RuleContainModelFields(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
  const (
    ruleName       = "contain-model-fields"
    defaultMessage = "Contain invalid fields: `:field`"
  )
  // g.Log().Debugf("%s: rule: %v value: %v message: %v data: %+v", ruleName, rule, value, message, data)

  // Set default message
  if message == "" {
    message = defaultMessage
  }

  // re := regexp.MustCompile(fmt.Sprintf(`^(?P<ruleName>%s):(?P<model>\w+),?(?P<modelFields>\w+,?)?$`, ruleName))
  // match := re.FindStringSubmatch(rule)

  var match []string
  splits := strings.Split(rule, ":")
  match = append(match, rule, splits[0])
  // append fields
  if len(splits) == 2 {
    fieldsStr := splits[1]
    splitFields := strings.Split(fieldsStr, ",")
    if len(splitFields) > 0 {
      match = append(match, splitFields...)
    }
  }

  g.Log().Debugf("result of %s resolution: %v", ruleName, match)

  var (
    modelName = ""
    // modelFields = match[3:]
    modelFields = garray.NewStrArray()
  )

  if len(match) > 2 {
    modelName = match[2]
    modelFields = garray.NewStrArrayFrom(match[3:])
  }

  // Check model fields by model name
  if modelFields.Len() == 0 {
    // Default set all model fields by model name if no specify fields.
    // Golang do not support dynamic load class so if model name do not registered, allow all fields
    if v, ok := modelFieldsMap[modelName]; ok {
      modelFields.SetArray(v)
    }
  }

  if v, ok := value.(string); ok {
    fields := strings.Split(v, ",")
    for _, field := range fields {
      if modelFields.Len() > 0 && !modelFields.Contains(field) {
        g.Log().Debugf("rule %s expect model fields: %v", ruleName, modelFields.Join(","))
        return errors.New(strings.Replace(message, ":field", field, -1))
      }
    }
  }

  return nil
}
