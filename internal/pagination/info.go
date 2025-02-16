// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package pagination

import (
	"math"
	"net/http"
	"strconv"
)

// maxPagesToShow is the maximum number of page links to show in front-end.
const maxPagesToShow = 5

// Info represents pagination information.
type Info struct {
	Pages        []int
	CurrentPage  int
	TotalPages   int
	ItemsPerPage int
	HasPrev      bool
	HasNext      bool
}

// New creates a new pagination Info based on the total number of items and
// current page.
func New(totalItems, currentPage, itemsPerPage int) Info {
	totalPages := int(math.Ceil(float64(totalItems) / float64(itemsPerPage)))
	if totalPages < 1 {
		totalPages = 1
	}

	if currentPage < 1 {
		currentPage = 1
	}

	if currentPage > totalPages {
		currentPage = totalPages
	}

	var (
		startPage = currentPage - (maxPagesToShow / 2)
		endPage   = currentPage + (maxPagesToShow / 2)
	)

	if startPage < 1 {
		endPage += (1 - startPage)

		startPage = 1
	}

	if endPage > totalPages {
		startPage -= (endPage - totalPages)
		if startPage < 1 {
			startPage = 1
		}

		endPage = totalPages
	}

	pages := make([]int, 0, endPage-startPage+1)
	for i := startPage; i <= endPage; i++ {
		pages = append(pages, i)
	}

	return Info{
		CurrentPage:  currentPage,
		TotalPages:   totalPages,
		HasPrev:      currentPage > 1,
		HasNext:      currentPage < totalPages,
		Pages:        pages,
		ItemsPerPage: itemsPerPage,
	}
}

// Offset returns the offset for database queries.
func (p Info) Offset() int {
	return (p.CurrentPage - 1) * p.ItemsPerPage
}

// Limit returns the limit for database queries.
func (p Info) Limit() int {
	return p.ItemsPerPage
}

// GetPage extracts the page number from the request.
func GetPage(r *http.Request) int {
	page := r.URL.Query().Get("page")
	if page == "" {
		return 1
	}

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		return 1
	}

	return pageNum
}
