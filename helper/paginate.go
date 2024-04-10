package helper

type Paginate interface {
	GetPaginateData(page, limit int) (int, int, int)
	GetTotalPages(totalData, limit int) int
}

type paginate struct{}

func NewPaginate() Paginate {
	return &paginate{}
}

func (p *paginate) GetPaginateData(page, limit int) (int, int, int) {
	if page <= 0 {
		page = 1
	}
	maxLimit := 100
	minLimit := 10

	switch {
	case limit > maxLimit:
		limit = maxLimit
	case limit <= 0:
		limit = minLimit
	}

	offset := (page - 1) * limit

	return page, offset, limit
}

func (p *paginate) GetTotalPages(totalData, limit int) int {
	if totalData <= 0 {
		return 1
	}

	totalPages := totalData / limit
	if totalData%limit > 0 {
		totalPages++
	}

	return totalPages
}
