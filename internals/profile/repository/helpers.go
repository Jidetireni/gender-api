package repository

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

var validSQLIdentifier = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*(\.[a-zA-Z_][a-zA-Z0-9_]*)?$`)

type QueryType int
type SortOrder string

const (
	QueryTypeSelect QueryType = iota
	QueryTypeCount

	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

type QueryOptions struct {
	Page  uint32
	Limit uint32
	Sort  *string
	Type  *QueryType
}

type SortResult struct {
	SortColumn string
	SortOrder  SortOrder
}

type ListResult[T any] struct {
	Items []*T
	Total int64
	Page  uint32
	Limit uint32
}

var allowedSortColumns = map[string]struct{}{
	"age":                {},
	"created_at":         {},
	"gender_probability": {},
}

func validateSQLIdentifier(identifier string) error {
	if !validSQLIdentifier.MatchString(identifier) {
		return fmt.Errorf("invalid SQL identifier: %q", identifier)
	}
	return nil
}

func parseSort(sort *string) (SortResult, error) {
	if sort == nil {
		return SortResult{
			SortColumn: "created_at",
			SortOrder:  SortOrderDesc,
		}, nil
	}

	parts := strings.Split(*sort, ":")
	if len(parts) != 2 {
		return SortResult{}, errors.New("invalid sort format")
	}

	sortColumn := parts[0]

	if err := validateSQLIdentifier(sortColumn); err != nil {
		return SortResult{}, err
	}

	if _, ok := allowedSortColumns[sortColumn]; !ok {
		return SortResult{}, fmt.Errorf("unsupported sort column: %q", sortColumn)
	}

	sortOrder := SortOrder(parts[1])
	switch sortOrder {
	case SortOrderAsc:
		return SortResult{
			SortColumn: sortColumn,
			SortOrder:  SortOrderAsc,
		}, nil
	case SortOrderDesc:
		return SortResult{
			SortColumn: sortColumn,
			SortOrder:  SortOrderDesc,
		}, nil
	default:
		return SortResult{}, errors.New("invalid sort order, must be asc or desc")
	}
}

func ApplyPagination(builder sq.SelectBuilder, queryOptions QueryOptions) (sq.SelectBuilder, error) {
	sortResult, err := parseSort(queryOptions.Sort)
	if err != nil {
		return builder, err
	}

	orderClause := fmt.Sprintf("%s %s", sortResult.SortColumn, sortResult.SortOrder)
	builder = builder.OrderBy(orderClause)

	limit := queryOptions.Limit
	if limit == 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	builder = builder.Limit(uint64(limit))

	page := queryOptions.Page
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit
	builder = builder.Offset(uint64(offset))

	return builder, nil
}
