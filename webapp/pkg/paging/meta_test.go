package paging

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewMetaInfo(t *testing.T) {
	tests := []struct {
		tag                           string
		offset, limit, total          int
		count                         bool
		expectedOffset, expectedLimit int
	}{
		{"t1", 1, 20, 50, true, 1, 20},
		{"t2", 2, 20, 50, true, 2, 20},
		{"t3", -2, 20, 50, false, DefaultOffset, 20},
		{"t4", 4, 20, 50, false, 4, 20},
		{"t5", 0, -100, 50, false, 0, DefaultLimit},
	}

	for _, test := range tests {
		p := NewMetaInfo(test.limit, test.offset, test.count, Unsorted)
		assert.Equal(t, test.expectedLimit, p.Limit, test.tag)
		assert.Equal(t, test.expectedOffset, p.Offset, test.tag)
	}
}

func Test_parseInt(t *testing.T) {
	type args struct {
		value        string
		defaultValue int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"t1", args{"123", 100}, 123},
		{"t2", args{"", 100}, 100},
		{"t3", args{"a", 100}, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseInt(tt.args.value, tt.args.defaultValue); got != tt.want {
				t.Errorf("parseInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_NewFromRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com?limit=2&offset=20", bytes.NewBufferString(""))
	p := NewMetaInfoFromRequest(req)
	assert.Equal(t, 2, p.Limit)
	assert.Equal(t, 20, p.Offset)
}
