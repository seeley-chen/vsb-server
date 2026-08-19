package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

var timeType = reflect.TypeOf(time.Time{})

// ValidationError 参数校验错误，Error() 可直接作为 API message 返回。
type ValidationError struct {
	Path string
	Msg  string
}

func (e *ValidationError) Error() string {
	return e.Msg
}

func requiredErr(path string) error {
	return &ValidationError{Path: path, Msg: path + " is required"}
}

func typeErr(path, want string) error {
	return &ValidationError{Path: path, Msg: path + " should be " + want}
}

func displayPath(path string) string {
	if path == "" {
		return "request"
	}
	return path
}

// BindAndValidate 解析 JSON 并按目标结构体的 json / validate / enums tag 校验。
// 空值：`{path} is required`；类型不符：`{path} should be {type}`；
// 对象数组会定位到具体下标和 key，如 `specs[0].price should be number`。
func BindAndValidate(body []byte, dst any) error {
	body = bytes.TrimSpace(body)
	if len(body) == 0 {
		return requiredErr("request body")
	}

	var raw any
	if err := json.Unmarshal(body, &raw); err != nil {
		return &ValidationError{Msg: "invalid request"}
	}

	if dst == nil {
		return &ValidationError{Msg: "invalid request"}
	}
	if err := validateJSON(raw, reflect.TypeOf(dst), "", false, nil); err != nil {
		return err
	}
	if err := json.Unmarshal(body, dst); err != nil {
		return &ValidationError{Msg: "invalid request"}
	}
	return ValidateStruct(dst)
}

// ValidateStruct 去掉字符串首尾空白后，按 required / enums 及嵌套路径校验。
func ValidateStruct(v any) error {
	if v == nil {
		return requiredErr("request")
	}
	return validateGoValue(reflect.ValueOf(v), "", false, nil)
}

func unwrapType(t reflect.Type) reflect.Type {
	for t != nil && t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

func validateJSON(raw any, typ reflect.Type, path string, required bool, enums []string) error {
	typ = unwrapType(typ)
	if typ == nil {
		return nil
	}

	if raw == nil {
		if required {
			return requiredErr(displayPath(path))
		}
		return nil
	}

	switch typ.Kind() {
	case reflect.String:
		s, ok := raw.(string)
		if !ok {
			return typeErr(displayPath(path), "string")
		}
		if required && strings.TrimSpace(s) == "" {
			return requiredErr(displayPath(path))
		}
		if len(enums) > 0 && strings.TrimSpace(s) != "" && !ContainsString(enums, s) {
			return typeErr(displayPath(path), joinOr(enums))
		}
		return nil

	case reflect.Bool:
		if _, ok := raw.(bool); !ok {
			return typeErr(displayPath(path), "boolean")
		}
		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, ok := asJSONNumber(raw)
		if !ok {
			return typeErr(displayPath(path), "integer")
		}
		if n != float64(int64(n)) {
			return typeErr(displayPath(path), "integer")
		}
		return nil

	case reflect.Float32, reflect.Float64:
		if _, ok := asJSONNumber(raw); !ok {
			return typeErr(displayPath(path), "number")
		}
		return nil

	case reflect.Slice, reflect.Array:
		arr, ok := raw.([]any)
		if !ok {
			return typeErr(displayPath(path), "array")
		}
		if required && len(arr) == 0 {
			return requiredErr(displayPath(path))
		}
		elemTyp := typ.Elem()
		for i, elem := range arr {
			elemPath := fmt.Sprintf("%s[%d]", displayPath(path), i)
			if err := validateJSON(elem, elemTyp, elemPath, true, nil); err != nil {
				return err
			}
		}
		return nil

	case reflect.Map:
		obj, ok := raw.(map[string]any)
		if !ok {
			return typeErr(displayPath(path), "object")
		}
		if required && len(obj) == 0 {
			return requiredErr(displayPath(path))
		}
		valTyp := typ.Elem()
		for k, val := range obj {
			p := k
			if path != "" {
				p = path + "." + k
			}
			if err := validateJSON(val, valTyp, p, false, nil); err != nil {
				return err
			}
		}
		if required && unwrapType(valTyp).Kind() == reflect.String && !jsonStringMapHasValue(obj) {
			return requiredErr(displayPath(path))
		}
		return nil

	case reflect.Struct:
		if typ == timeType {
			if _, ok := raw.(string); !ok {
				return typeErr(displayPath(path), "string")
			}
			return nil
		}
		obj, ok := raw.(map[string]any)
		if !ok {
			return typeErr(displayPath(path), "object")
		}
		return validateJSONStruct(obj, typ, path)

	default:
		return nil
	}
}

func validateJSONStruct(obj map[string]any, typ reflect.Type, prefix string) error {
	typ = unwrapType(typ)
	if typ == nil || typ.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if !f.IsExported() {
			continue
		}
		tag := f.Tag.Get("json")
		if tag == "-" {
			continue
		}
		if f.Anonymous && (tag == "" || strings.HasPrefix(tag, ",")) {
			if err := validateJSONStruct(obj, f.Type, prefix); err != nil {
				return err
			}
			continue
		}

		name := jsonName(f)
		path := name
		if prefix != "" {
			path = prefix + "." + name
		}
		val, ok := obj[name]
		if !ok {
			val = nil
		}
		if err := validateJSON(val, f.Type, path, hasRequired(f.Tag.Get("validate")), parseCSV(f.Tag.Get("enums"))); err != nil {
			return err
		}
	}
	return nil
}

func validateGoValue(v reflect.Value, path string, required bool, enums []string) error {
	if !v.IsValid() {
		if required {
			return requiredErr(displayPath(path))
		}
		return nil
	}
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			if required {
				return requiredErr(displayPath(path))
			}
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		s := strings.TrimSpace(v.String())
		if v.CanSet() {
			v.SetString(s)
		}
		if required && s == "" {
			return requiredErr(displayPath(path))
		}
		if s != "" && len(enums) > 0 && !ContainsString(enums, s) {
			return typeErr(displayPath(path), joinOr(enums))
		}
		return nil

	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return nil

	case reflect.Slice, reflect.Array:
		if required && v.Len() == 0 {
			return requiredErr(displayPath(path))
		}
		for i := 0; i < v.Len(); i++ {
			elemPath := fmt.Sprintf("%s[%d]", displayPath(path), i)
			if err := validateGoValue(v.Index(i), elemPath, true, nil); err != nil {
				return err
			}
		}
		return nil

	case reflect.Map:
		if v.IsNil() {
			if required {
				return requiredErr(displayPath(path))
			}
			return nil
		}
		if v.CanSet() && v.Type().Key().Kind() == reflect.String && v.Type().Elem().Kind() == reflect.String {
			trimmed := reflect.MakeMapWithSize(v.Type(), v.Len())
			iter := v.MapRange()
			for iter.Next() {
				trimmed.SetMapIndex(iter.Key(), reflect.ValueOf(strings.TrimSpace(iter.Value().String())))
			}
			v.Set(trimmed)
		}
		if required && !goMapHasValue(v) {
			return requiredErr(displayPath(path))
		}
		iter := v.MapRange()
		for iter.Next() {
			k := fmt.Sprint(iter.Key().Interface())
			p := k
			if path != "" {
				p = path + "." + k
			}
			if err := validateGoValue(iter.Value(), p, false, nil); err != nil {
				return err
			}
		}
		return nil

	case reflect.Struct:
		if v.Type() == timeType {
			if required && v.IsZero() {
				return requiredErr(displayPath(path))
			}
			return nil
		}
		return validateGoStruct(v, path)

	default:
		if required && v.IsZero() {
			return requiredErr(displayPath(path))
		}
		return nil
	}
}

func validateGoStruct(v reflect.Value, prefix string) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		tag := f.Tag.Get("json")
		if tag == "-" {
			continue
		}
		fv := v.Field(i)
		if f.Anonymous && (tag == "" || strings.HasPrefix(tag, ",")) {
			if err := validateGoStruct(fv, prefix); err != nil {
				return err
			}
			continue
		}
		name := jsonName(f)
		path := name
		if prefix != "" {
			path = prefix + "." + name
		}
		if err := validateGoValue(fv, path, hasRequired(f.Tag.Get("validate")), parseCSV(f.Tag.Get("enums"))); err != nil {
			return err
		}
	}
	return nil
}

func jsonName(f reflect.StructField) string {
	tag := f.Tag.Get("json")
	if tag == "" || tag == "-" {
		return f.Name
	}
	name := strings.Split(tag, ",")[0]
	if name == "" {
		return f.Name
	}
	return name
}

func hasRequired(tag string) bool {
	for _, p := range strings.Split(tag, ",") {
		if strings.TrimSpace(p) == "required" {
			return true
		}
	}
	return false
}

func parseCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func joinOr(vals []string) string {
	return strings.Join(vals, " or ")
}

func jsonStringMapHasValue(obj map[string]any) bool {
	for _, val := range obj {
		s, ok := val.(string)
		if ok && strings.TrimSpace(s) != "" {
			return true
		}
	}
	return false
}

func goMapHasValue(v reflect.Value) bool {
	if v.Len() == 0 {
		return false
	}
	if v.Type().Elem().Kind() != reflect.String {
		return true
	}
	iter := v.MapRange()
	for iter.Next() {
		if strings.TrimSpace(iter.Value().String()) != "" {
			return true
		}
	}
	return false
}

func asJSONNumber(raw any) (float64, bool) {
	switch n := raw.(type) {
	case float64:
		return n, true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	default:
		return 0, false
	}
}
