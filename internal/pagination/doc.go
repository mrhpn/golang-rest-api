// Package pagination provides reusable utilities for pagination, searching, and sorting
// in list endpoints across the application.
//
// Features:
//   - Pagination parameter normalization
//   - Pagination metadata building
//   - Search across multiple columns (case-insensitive, with optional exact match)
//   - Field validation for sorting (supports snake_case and camelCase)
//   - GORM query helpers for pagination, search, and sorting
//   - Struct-based QueryOptions to prevent parameter order mistakes
//
// Example usage:
//
//	// In handler
//	var query pagination.ListQuery
//	if err := c.ShouldBindQuery(&query); err != nil {
//	    httpx.FailWithError(c, err)
//	    return
//	}
//
//	// Define model-specific settings
//	allowedSortFields := []string{"email", "role", "created_at"}
//	searchableColumns := []string{"email", "role"}
//
//	// Create QueryOptions struct (prevents parameter order mistakes)
//	opts := pagination.NewQueryOptions(&query, &User{}, allowedSortFields, searchableColumns)
//
//	// In service - pass single struct instead of many parameters
//	users, meta, err := userService.List(ctx, opts)
//
//	// In repository - use QueryOptions struct
//	baseQuery := db.Model(&User{})
//	countQuery := pagination.CountQuery(baseQuery, opts)
//	dataQuery, _ := pagination.PaginateQuery(baseQuery, opts)
//
//	// Build pagination metadata
//	meta := pagination.BuildPaginationMeta(opts.Page, opts.Limit, int(total))
package pagination
