package paging

import (
	"net/http"
	"strconv"
)

const (
	DefaultLimit    = 100
	DefaultOffset   = 0
	OffsetParamName = "offset"
	LimitParamName  = "limit"
	CountParamName  = "count"
)

type MetaInfo struct {
	Offset int  `json:"offset" form:"offset"`
	Limit  int  `json:"limit" form:"limit"`
	Total  int  `json:"total"`
	Count  bool `json:"count" form:"count"`
	Sort   Sort `json:"sort" `
}

func NewMetaInfo(limit, offset int, count bool, sort Sort) *MetaInfo {
	if limit <= 0 {
		limit = DefaultLimit
	}

	if offset < 0 {
		offset = DefaultOffset
	}

	return &MetaInfo{offset, limit, -1, count, sort}
}

func NewMetaInfoFromRequest(request *http.Request) MetaInfo {
	limit := parseInt(request.URL.Query().Get(LimitParamName), DefaultLimit)
	offset := parseInt(request.URL.Query().Get(OffsetParamName), DefaultOffset)
	count := parseBool(request.URL.Query().Get(CountParamName), true)
	sort := NewSortFromRequest(request)
	return *NewMetaInfo(limit, offset, count, sort)
}

func parseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	if result, err := strconv.Atoi(value); err == nil {
		return result
	}
	return defaultValue
}

func parseBool(value string, defaultValue bool) bool {
	if value == "" {
		return true
	}
	if result, err := strconv.ParseBool(value); err == nil {
		return result
	}
	return defaultValue
}
