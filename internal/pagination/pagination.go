package pagination

import (
	"strings"

	"gorm.io/gorm"

	"github.com/mrhpn/go-rest-api/internal/constants"
	"github.com/mrhpn/go-rest-api/internal/httpx"
	"github.com/mrhpn/go-rest-api/internal/stringx"
)

// SortSearchPolicy contains SortableCols and SearchableCols. Contents containing in both must match DB column names or model field names
// Mapping assumes snake_case DB columns
type SortSearchPolicy struct {
	SortableCols   []string
	SearchableCols []string
}

// QueryList represents query parameters for list endpoints
type QueryList struct {
	Page          int      `form:"page" binding:"omitempty,min=1"`
	Limit         int      `form:"limit" binding:"omitempty,min=1"`
	Search        string   `form:"search" binding:"omitempty"`
	SearchColumns []string `form:"search_columns" binding:"omitempty"`
	SortBy        string   `form:"sort_by" binding:"omitempty"`
	Order         string   `form:"order" binding:"omitempty,oneof=asc desc ASC DESC"`
	ExactMatch    bool     `form:"exact_match" binding:"omitempty"`
}

// QueryOptions encapsulates all parameters for pagination, search, and sorting queries
// This prevents parameter order mistakes and makes function signatures cleaner
type QueryOptions struct {
	// Pagination
	Page   int
	Limit  int
	Offset int // Calculated from Page and Limit

	// Search
	Search            string
	SearchColumns     []string
	SearchableColumns []string // Default searchable columns for the model
	ExactMatch        bool

	// Sorting
	SortBy          string
	Order           string
	SortableColumns map[string]string // Fields allowed for sorting

	// Model for field validation
	Model any
}

// NewQueryOptions creates a QueryOptions from QueryList and model-specific settings
func NewQueryOptions(ql *QueryList, sortSearchPolicy SortSearchPolicy) *QueryOptions {
	// Ensure pagination params
	if ql.Page < 1 {
		ql.Page = constants.PaginationDefaultPage
	}
	if ql.Limit < 1 {
		ql.Limit = constants.PaginationDefaultLimit
	}
	if ql.Limit > constants.PaginationMaxLimit {
		ql.Limit = constants.PaginationMaxLimit
	}

	return &QueryOptions{
		Page:              ql.Page,
		Limit:             ql.Limit,
		Offset:            (ql.Page - 1) * ql.Limit,
		Search:            ql.Search,
		SearchColumns:     ql.SearchColumns,
		SearchableColumns: sortSearchPolicy.SearchableCols,
		ExactMatch:        ql.ExactMatch,
		SortBy:            ql.SortBy,
		Order:             normalizeOrder(ql.Order),
		SortableColumns:   buildSortableFieldMap(sortSearchPolicy.SortableCols),
	}
}

// Paginate is a GORM Scope for data fetching
func Paginate(opts *QueryOptions) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 1. Apply Search
		db = applySearch(db, opts)

		// 2. Apply Sorting
		sortField := "created_at" // default
		if opts.SortBy != "" {
			if col, ok := opts.SortableColumns[stringx.ToSnakeCase(opts.SortBy)]; ok {
				sortField = col
			}
		}
		db = db.Order(sortField + " " + opts.Order) // opts.Order is safe. already validated in NewQueryOptions!

		// 3. Apply Pagination
		return db.Limit(opts.Limit).Offset(opts.Offset)
	}
}

// SearchScope is a GORM Scope for the count query
func SearchScope(opts *QueryOptions) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return applySearch(db, opts)
	}
}

// applySearch applies search conditions to the query
func applySearch(db *gorm.DB, opts *QueryOptions) *gorm.DB {
	if opts.Search == "" {
		return db
	}

	cols := opts.SearchColumns
	if len(cols) == 0 {
		cols = opts.SearchableColumns
	}
	if len(cols) == 0 {
		return db
	}

	var conditions []string
	var args []any
	for _, col := range cols {
		columnName := stringx.ToSnakeCase(col)
		if opts.ExactMatch {
			conditions = append(conditions, columnName+" = ? ")
			args = append(args, opts.Search)
		} else {
			conditions = append(conditions, columnName+" ILIKE ?")
			args = append(args, "%"+opts.Search+"%")
		}
	}
	return db.Where(strings.Join(conditions, " OR "), args...)
}

// BuildMeta creates pagination metadata from page, limit, and total count
// If page is out of bounds, has_prev is true if there are valid previous pages, and has_next is always false for out-of-bounds pages
func BuildMeta(opts *QueryOptions, total int64) *httpx.PaginationMeta {
	totalInt := int(total)
	totalPages := (totalInt + opts.Limit - 1) / opts.Limit
	if totalPages == 0 {
		totalPages = 1
	}

	isValidPage := opts.Page >= 1 && opts.Page <= totalPages

	return &httpx.PaginationMeta{
		Page:       opts.Page,
		Limit:      opts.Limit,
		Total:      totalInt,
		TotalPages: totalPages,
		HasNext:    isValidPage && opts.Page < totalPages,
		HasPrev:    isValidPage && opts.Page > 1,
	}
}

func buildSortableFieldMap(fields []string) map[string]string {
	m := make(map[string]string, len(fields))
	for _, f := range fields {
		key := stringx.ToSnakeCase(f)
		m[key] = key
	}
	return m
}

func normalizeOrder(order string) string {
	switch strings.ToUpper(order) {
	case "ASC":
		return "ASC"
	default:
		return "DESC"
	}
}
