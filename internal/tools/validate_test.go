package tools

import (
	"testing"
)

type specItem struct {
	Price  float64 `json:"price" validate:"required"`
	Size   string  `json:"size" validate:"required"`
	Status bool    `json:"status"`
}

type sampleReq struct {
	Username string     `json:"username" validate:"required"`
	Identity string     `json:"identity" validate:"required" enums:"admin,user"`
	Specs    []specItem `json:"specs"`
}

func TestBindAndValidate(t *testing.T) {
	cases := []struct {
		name string
		body string
		want string
	}{
		{"missing required", `{}`, "username is required"},
		{"wrong primitive type", `{"username":1}`, "username should be string"},
		{"enum", `{"username":"a","identity":"x"}`, "identity should be admin or user"},
		{"nested type", `{"username":"a","identity":"admin","specs":[{"price":"10","size":"M"}]}`, "specs[0].price should be number"},
		{"nested required", `{"username":"a","identity":"admin","specs":[{"price":1}]}`, "specs[0].size is required"},
		{"nested object type", `{"username":"a","identity":"admin","specs":["x"]}`, "specs[0] should be object"},
		{"array type", `{"username":"a","identity":"admin","specs":{}}`, "specs should be array"},
		{"whitespace required", `{"username":"  ","identity":"admin"}`, "username is required"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var req sampleReq
			err := BindAndValidate([]byte(tc.body), &req)
			if err == nil {
				t.Fatalf("want %q, got nil", tc.want)
			}
			if err.Error() != tc.want {
				t.Fatalf("want %q, got %q", tc.want, err.Error())
			}
		})
	}

	var req sampleReq
	if err := BindAndValidate([]byte(`{"username":"a","identity":"admin","specs":[{"price":1.2,"size":"M","status":false}]}`), &req); err != nil {
		t.Fatal(err)
	}
	if req.Username != "a" {
		t.Fatalf("username not trimmed/bound: %q", req.Username)
	}
}

type localeReq struct {
	Name map[string]string `json:"name" validate:"required"`
}

type permItem struct {
	Path     string     `json:"path" validate:"required"`
	Type     string     `json:"type" validate:"required" enums:"read,write"`
	Children []permItem `json:"children"`
}

type roleReq struct {
	Name        string     `json:"name" validate:"required"`
	Permissions []permItem `json:"permissions"`
}

func TestBindAndValidateNested(t *testing.T) {
	var loc localeReq
	if err := BindAndValidate([]byte(`{"name":{"zh-cn":"  "}}`), &loc); err == nil || err.Error() != "name is required" {
		t.Fatalf("locale empty: %v", err)
	}

	var role roleReq
	err := BindAndValidate([]byte(`{"name":"admin","permissions":[{"path":"/a","type":"x"}]}`), &role)
	if err == nil || err.Error() != "permissions[0].type should be read or write" {
		t.Fatalf("perm enum: %v", err)
	}

	err = BindAndValidate([]byte(`{"name":"admin","permissions":[{"path":"/a","type":"read","children":[{"type":"write"}]}]}`), &role)
	if err == nil || err.Error() != "permissions[0].children[0].path is required" {
		t.Fatalf("nested child: %v", err)
	}
}
