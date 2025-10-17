package dto

type PaginationQuery struct {
	Page       *int64  `query:"page" validate:"omitempty,gte=0"`
	Limit      *int64  `query:"limit" validate:"omitempty,gte=0"`
	NextCursor *string `query:"cursor"`
}

func (p *PaginationQuery) UseCursor() bool {
	return p.NextCursor != nil && *p.NextCursor != ""
}

func (p *PaginationQuery) SetDefaults() {
	if p.Limit == nil {
		defaultLimit := int64(10)
		p.Limit = &defaultLimit
	}
	if p.Page == nil {
		defaultPage := int64(1)
		p.Page = &defaultPage
	}
}
