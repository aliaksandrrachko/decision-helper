package paging

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewSortFromRequest(t *testing.T) {
	tests := []struct {
		tag          string
		url          string
		expectedSort Sort
	}{
		{"t1", "http://example.com?limit=2&offset=20&sort=user", NewSortBy("user")},
		{"t2", "http://example.com?limit=2&offset=20&sort=user.id&sort=user.name,ASC", NewSortBy("user.id", "user.name")},
		{"t3", "http://example.com?limit=2&offset=20&sort=user.id,DESC&sort=user.name,DESC", NewSortBy("user.id", "user.name").Descending()},
	}

	for _, test := range tests {
		req, _ := http.NewRequest(http.MethodGet, test.url, bytes.NewBufferString(""))
		sort := NewSortFromRequest(req)
		assert.Equal(t, test.expectedSort, sort, test.tag)
	}
}

func Test_parseOrder(t *testing.T) {
	tests := []struct {
		tag                  string
		value                string
		expectedDirection    Direction
		expectedPropertyName string
		isErr                bool
	}{
		{"t1", "id,asc", Ascending, "id", false},
		{"t2", "name", Ascending, "name", false},
		{"t3", "user.name,DESC", Descending, "user.name", false},
		{"t4", "user.type.id,asc", Ascending, "user.type.id", false},
		{"t5", "user.id,desc", Descending, "user.id", false},
	}

	for _, test := range tests {
		p, err := parseOrder(test.value)
		assert.Equal(t, test.isErr, err != nil, test.tag)
		assert.Equal(t, test.expectedDirection, p.Direction, test.tag)
		assert.Equal(t, test.expectedPropertyName, p.Property, test.tag)
	}
}
