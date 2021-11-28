package sorm

import (
  "context"
  "fmt"
  "testing"
)

type User struct {
  Id       uint64 `orm:"id,primary" json:"id"       description:""`                  //
  Name     string `orm:"name"       json:"name"     description:"user name"`         // user name
  Nick     string `orm:"nick"       json:"nick"     description:"user nickname"`     // user nickname
  Password string `orm:"password"   json:"password" description:"user password"`     // user password
  Email    string `orm:"email"      json:"email"    description:"user email"`        // user email
  Age      uint   `orm:"age"        json:"age"      description:"user age"`          // user age
  HeadImg  string `orm:"head_img"   json:"headImg"  description:"user head img url"` // user head img url
}

func TestRegisterModelColumnsMap(t *testing.T) {
  type args struct {
    m map[string]interface{}
  }
  tests := []struct {
    name string
    args args
  }{
    {"register model columns", args{map[string]interface{}{"user": User{}}}},
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      RegisterModelColumnsMap(tt.args.m)
    })
  }
}

func Test_checkOrderColumns(t *testing.T) {
  // register model columns
  RegisterModelColumnsMap(map[string]interface{}{"user": User{}})

  type args struct {
    model   string
    columns []string
    value   string
    message string
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    // Success cases.
    {"no rule value", args{"", []string{}, "id,name", ""}, false},
    {"no model", args{"", []string{"id", "name", "age"}, "id,name", ""}, false},
    {"no model columns(model has registered by RegisterModelColumnsMap)", args{"user", []string{}, "id,name", ""}, false},
    {"have model columns", args{"user", []string{"id", "name", "age"}, "id,name", ""}, false},
    //  Failure cases.
    {"contain invalid model columns", args{"user", []string{}, "id,invalidCol", ""}, true},
    {"case-sensitive", args{"user", []string{"id", "name"}, "Id", ""}, true},
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if err := checkOrderColumns(tt.args.model, tt.args.columns, tt.args.value,
        tt.args.message); (err != nil) != tt.wantErr {
        t.Errorf("RuleContainModelFields() error = %v, wantErr %v", err, tt.wantErr)
      } else if err != nil {
        t.Log("wantErr: ", err)
      }
    })
  }
}

func TestRuleOrderBy(t *testing.T) {
  // register model columns
  RegisterModelColumnsMap(map[string]interface{}{"user": User{}})

  var (
    ctx      context.Context = context.Background()
    ruleName                 = OrderByRuleName
    genRule                  = func(s string) string {
      if s == "" {
        return ruleName
      }
      return fmt.Sprintf("%s:%s", ruleName, s)
    }
  )
  type args struct {
    ctx     context.Context
    rule    string
    value   interface{}
    message string
    data    interface{}
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    // Success cases.
    {"no rule", args{ctx, genRule(""), nil, "", map[string]interface{}{"orderColumns": "id,name"}}, false},
    {"no model", args{ctx, genRule(""), nil, "", map[string]interface{}{"orderColumns": "id,name"}}, false},
    {"no model columns(model has registered by RegisterModelColumnsMap)", args{ctx, genRule("user"), nil, "", map[string]interface{}{"orderColumns": "id,name"}}, false},
    {"have model columns", args{ctx, genRule("user,id,name,age"), nil, "",
      map[string]interface{}{"orderColumns": "id," +
          "name"}}, false},
    // Failure cases.
    {"no model column and contain invalid columns", args{ctx, genRule("user"), nil, "",
      map[string]interface{}{"orderColumns": "id,invalidcol"}},
      true},
    {"contain invalid columns", args{ctx, genRule("user,id,name"), nil, "",
      map[string]interface{}{"orderColumns": "id,name,age"}}, true},
    {"case-sensitive", args{ctx, genRule("user,id,name"), nil, "", map[string]interface{}{"orderColumns": "id,Name"}},
      true},
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if err := RuleOrderBy(tt.args.ctx, tt.args.rule, tt.args.value, tt.args.message,
        tt.args.data); (err != nil) != tt.wantErr {
        t.Errorf("RuleContainModelFields() error = %v, wantErr %v", err, tt.wantErr)
      } else if err != nil {
        t.Log("wantErr: ", err)
      }
    })
  }
}
