package database

import "log"

// PaginatedResult represents a generic paginated result
type PaginatedResult[T any] struct {
	Items    []T
	NextPage int
}

// PaginationFetcher is a function type that fetches paginated data
type PaginationFetcher[T any] func(page, limit int) (PaginatedResult[T], error)

// FetchAll retrieves all items using pagination
func FetchAll[T any](fetcher PaginationFetcher[T], initialPage, pageSize int) ([]T, error) {
	var allItems []T
	currentPage := initialPage

	for {
		result, err := fetcher(currentPage, pageSize)
		if err != nil {
			log.Printf("Failed to fetch items: %v", err)
			return nil, err
		}

		allItems = append(allItems, result.Items...)

		if result.NextPage == currentPage {
			break
		}

		currentPage = result.NextPage
	}

	return allItems, nil
}
