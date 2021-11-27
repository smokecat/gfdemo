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

func TestRegisterModelFieldsMap(t *testing.T) {
	type args struct {
		m map[string]interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"register model fields", args{map[string]interface{}{"user": User{}}}},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterModelFieldsMap(tt.args.m)
		})
	}
}

func TestRuleContainModelFields(t *testing.T) {
	// register model fields
	RegisterModelFieldsMap(map[string]interface{}{"user": User{}})

	var (
		ctx        context.Context = context.Background()
		rulePrefix                 = "contain-model-fields"
		genRule                    = func(s string) string {
			if s == "" {
				return rulePrefix
			}
			return fmt.Sprintf("%s:%s", rulePrefix, s)
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
		{"empty rule", args{ctx: ctx, rule: genRule(""), value: "Id", message: "", data: "data"}, false},
		{"not registered model", args{ctx: ctx, rule: genRule("mdl,Id"), value: "Id", message: "", data: "data"}, false},
		{"rule contains only model", args{ctx: ctx, rule: genRule("user"), value: "Id", message: "",
			data: "data"}, false},
		{"rule contains fields", args{ctx: ctx, rule: genRule("user,Id"), value: "Id", message: "", data: "data"},
			false},
		{"pass valid fields", args{ctx: ctx, rule: genRule("user,Id,Name,Age"), value: "Id,Name",
			message: "", data: "data"}, false},
		{"pass invalid fields", args{ctx: ctx, rule: genRule("user,Id"), value: "Id,Name", message: "",
			data: "data"}, true},
		{"pass invalid fields with space", args{ctx: ctx, rule: genRule("user,Id,Name"), value: "Id, Name", message: "",
			data: "data"}, true},
		{"field name case-sensitive", args{ctx: ctx, rule: genRule("user,Id"), value: "id", message: "",
			data: "data"}, true},
		{"custom message", args{ctx: ctx, rule: genRule("user,Id"), value: "id",
			message: "Custom message: contain invalid field: :field", data: "data"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RuleContainModelFields(tt.args.ctx, tt.args.rule, tt.args.value, tt.args.message, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("RuleContainModelFields() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Log("wantErr: ", err)
			}
		})
	}
}
