package helpers

import "ChainServer/internal/common/response"

func BuildPaginationMeta(limit, page, total int64, nextCursor any) *response.PaginationMeta {
	return &response.PaginationMeta{
		Limit:       int(limit),
		CurrentPage: int(page),
		Total:       int(total),
		NextCursor:  nextCursor,
		HasMore:     total > page*limit,
	}
}
