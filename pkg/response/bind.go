package response

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var validate = validator.New()

// BindJSON 解析 JSON body 并校验 validate tag，失败时写入 400 并返回 false。
func BindJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		Fail(w, http.StatusBadRequest, CodeBadRequest, "invalid request")
		return false
	}
	if err := validate.Struct(dst); err != nil {
		Fail(w, http.StatusBadRequest, CodeBadRequest, "invalid request: "+err.Error())
		return false
	}
	return true
}

// PathVar 取路径参数，为空时写入 400 并返回 false。
func PathVar(w http.ResponseWriter, r *http.Request, key string, errMsg ...string) (string, bool) {
	val := mux.Vars(r)[key]
	if val == "" {
		msg := key + " is required"
		if len(errMsg) > 0 && errMsg[0] != "" {
			msg = errMsg[0]
		}
		Fail(w, http.StatusBadRequest, CodeBadRequest, msg)
		return "", false
	}
	return val, true
}

// PageQuery 解析 pageIndex/pageSize，默认 1/20，pageSize 上限 100。
func PageQuery(r *http.Request) (pageIndex, pageSize int) {
	pageIndex, pageSize = 1, 20
	if p, err := strconv.Atoi(r.URL.Query().Get("pageIndex")); err == nil && p > 0 {
		pageIndex = p
	}
	if ps, err := strconv.Atoi(r.URL.Query().Get("pageSize")); err == nil && ps > 0 && ps <= 100 {
		pageSize = ps
	}
	return pageIndex, pageSize
}
