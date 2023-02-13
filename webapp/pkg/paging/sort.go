package paging

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var (
	Unsorted Sort = Sort{make([]Order, 0)}
)

type Sort struct {
	Orders []Order `json:"order"`
}

func NewSortBy(properties ...string) Sort {
	sort := &Sort{make([]Order, 0)}
	for _, p := range properties {
		sort.Orders = append(sort.Orders, NewOrderAsc(p))
	}
	return *sort
}

func NewSortFromRequest(request *http.Request) Sort {
	return parseSort(request)
}

func (s Sort) Ascending() Sort {
	return s.withDirection(Ascending)
}

func (s Sort) Descending() Sort {
	return s.withDirection(Descending)
}

func (s Sort) withDirection(direction Direction) Sort {
	newSort := &Sort{make([]Order, 0)}
	for _, o := range s.Orders {
		newSort.Orders = append(newSort.Orders, Order{direction, o.Property})
	}
	return *newSort
}

func parseSort(request *http.Request) Sort {
	v := request.URL.Query() // Values map[string][]string

	if v == nil {
		return Unsorted
	}

	orders := make([]Order, 0)
	for key, element := range v {
		if strings.EqualFold(key, SortParameterName) {
			for _, v := range element {
				if o, err := parseOrder(v); err == nil {
					orders = append(orders, o)
				}
			}
		}
	}

	return Sort{orders}
}

type Order struct {
	Direction Direction `json:"direction"`
	Property  string    `json:"property"`
}

func NewOrderAsc(property string) Order {
	return Order{Ascending, property}
}

func NewOrderDesc(property string) Order {
	return Order{Descending, property}
}

func NewOrderBy(property string) Order {
	return Order{DefaultDirection, property}
}

type Direction string

const (
	Ascending  Direction = "ASC"
	Descending Direction = "DESC"

	SortParameterName = "sort"
	DefaultDirection  = Ascending
)

func (d Direction) String() string {
	switch d {
	case Ascending:
		return "ASC"
	case Descending:
		return "DESC"
	default:
		return "UNKNOWN"
	}
}

// Parse parameters from string.
// Sorting criteria in the format: property(,asc|desc). Default sort order is ascending.
// Multiple sort criteria are allowed.
func parseOrder(value string) (Order, error) {
	regex := regexp.MustCompile("^([a-zA-Z_.]+)(,(?i)(asc|desc))?$")

	parseError := fmt.Errorf("Could not parse Order from string '%v' by pattern %v", value, regex)

	if ok := regex.MatchString(value); ok {
		v := regex.FindStringSubmatch(value)
		if len(v) < 1 {
			return Order{}, parseError
		}

		propertyName := v[1]
		if len(v) < 3 {
			return NewOrderBy(propertyName), nil
		}

		if strings.EqualFold(v[3], Descending.String()) {
			return NewOrderDesc(propertyName), nil
		}
		return NewOrderBy(propertyName), nil
	}

	return Order{}, parseError
}
